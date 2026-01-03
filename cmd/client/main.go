package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/example/hello-ct/gen/go/hello"
	"github.com/example/hello-ct/gen/go/hello/helloconnect"
)

func main() {
	baseURL := "http://localhost:8080"
	if url := os.Getenv("SERVER_URL"); url != "" {
		baseURL = url
	}

	// Create Connect client, wrap as Kleisli service
	connectClient := helloconnect.NewHelloServiceClient(http.DefaultClient, baseURL)
	client := hello.NewHelloServiceKFromConnectWithDefaults(connectClient, 5*time.Second, 3)

	ctx := context.Background()

	// 1. Simple call
	fmt.Println("=== Simple Call ===")
	resp, err := client.SayHello(&hello.SayHelloRequest{Name: "World"})(ctx)
	if err != nil {
		log.Fatalf("SayHello failed: %v", err)
	}
	fmt.Printf("Response: %s (count: %d)\n\n", resp.Greeting.Message, resp.Greeting.Count)

	// 2. Parallel calls with Traverse
	fmt.Println("=== Parallel Traverse ===")
	names := []string{"Alice", "Bob", "Charlie"}
	responses, err := hello.TraverseParallelNetworkOp(names, func(name string) hello.NetworkOp[*hello.SayHelloResponse] {
		return client.SayHello(&hello.SayHelloRequest{Name: name})
	})(ctx)
	if err != nil {
		log.Fatalf("Traverse failed: %v", err)
	}
	for _, r := range responses {
		fmt.Printf("  %s\n", r.Greeting.Message)
	}
	fmt.Println()

	// 3. Kleisli composition
	fmt.Println("=== Kleisli Bind ===")
	composed := hello.BindNetworkOp(
		client.SayHello(&hello.SayHelloRequest{Name: "Composed"}),
		func(resp *hello.SayHelloResponse) hello.NetworkOp[*hello.ListGreetingsResponse] {
			fmt.Printf("  First: %s\n", resp.Greeting.Message)
			return client.ListGreetings(&hello.ListGreetingsRequest{Limit: 5})
		},
	)
	listResp, err := composed(ctx)
	if err != nil {
		log.Fatalf("Composed failed: %v", err)
	}
	fmt.Printf("  Listed %d greetings\n\n", len(listResp.Greetings))

	// 4. Monoid combine
	fmt.Println("=== Monoid Combine ===")
	g1 := &hello.Greeting{Id: "1", Message: "Hello", Count: 5}
	g2 := &hello.Greeting{Id: "2", Message: "World", Count: 3}
	combined := hello.GreetingMonoid.Combine(g1, g2)
	fmt.Printf("  Combined count: %d (5+3)\n\n", combined.Count)

	// 5. FoldMap
	fmt.Println("=== FoldMap ===")
	greetings := []*hello.Greeting{{Count: 10}, {Count: 20}, {Count: 30}}
	total := hello.FoldMap(greetings, hello.Id[*hello.Greeting](), hello.GreetingMonoid)
	fmt.Printf("  Total count: %d\n\n", total.Count)

	fmt.Println("Done!")
}
