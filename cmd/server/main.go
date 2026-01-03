package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/example/hello-ct/gen/go/hello"
	"github.com/example/hello-ct/gen/go/hello/helloconnect"
	"github.com/example/hello-ct/internal/service"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Connect-Protocol-Version")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	ctx := context.Background()

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		projectID = "your-project-id"
	}

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	svc := service.NewHelloService(client)

	// Add logging middleware
	wrapped := hello.ApplyHelloServiceMiddleware(svc, &hello.HelloServiceMiddleware{
		SayHello:      loggingMiddleware[hello.SayHelloRequest, hello.SayHelloResponse]("SayHello"),
		ListGreetings: loggingMiddleware[hello.ListGreetingsRequest, hello.ListGreetingsResponse]("ListGreetings"),
	})

	// Create Connect handler
	handler := hello.NewHelloServiceConnectHandler(wrapped)

	mux := http.NewServeMux()
	path, h := helloconnect.NewHelloServiceHandler(handler)
	mux.Handle(path, h)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, h2c.NewHandler(corsMiddleware(mux), &http2.Server{})); err != nil {
		log.Fatal(err)
	}
}

func loggingMiddleware[Req, Resp any](method string) func(func(*Req) hello.NetworkOp[*Resp]) func(*Req) hello.NetworkOp[*Resp] {
	return func(next func(*Req) hello.NetworkOp[*Resp]) func(*Req) hello.NetworkOp[*Resp] {
		return func(req *Req) hello.NetworkOp[*Resp] {
			return func(ctx context.Context) (*Resp, error) {
				start := time.Now()
				resp, err := next(req)(ctx)
				log.Printf("[%s] duration=%v err=%v", method, time.Since(start), err)
				return resp, err
			}
		}
	}
}
