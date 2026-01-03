package service

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/example/hello-ct/gen/go/hello"
	"github.com/google/uuid"
)

type helloService struct {
	greetings *hello.GreetingCollection
}

// NewHelloService creates a HelloServiceK backed by Firestore
func NewHelloService(fsClient *firestore.Client) *hello.HelloServiceK {
	svc := &helloService{
		greetings: hello.NewGreetingCollection(fsClient),
	}

	return &hello.HelloServiceK{
		SayHello:      svc.sayHello,
		ListGreetings: svc.listGreetings,
	}
}

func (s *helloService) sayHello(req *hello.SayHelloRequest) hello.NetworkOp[*hello.SayHelloResponse] {
	return func(ctx context.Context) (*hello.SayHelloResponse, error) {
		id := uuid.New().String()

		greeting := &hello.Greeting{
			Id:      id,
			Message: "Hello, " + req.Name + "!",
			Count:   1,
		}

		// Save to Firestore - LiftDBOpToNetworkOp bridges the effect types
		saveOp := hello.LiftDBOpToNetworkOp(s.greetings.Doc(id).Set(greeting))
		if _, err := saveOp(ctx); err != nil {
			return nil, err
		}

		return &hello.SayHelloResponse{Greeting: greeting}, nil
	}
}

func (s *helloService) listGreetings(req *hello.ListGreetingsRequest) hello.NetworkOp[*hello.ListGreetingsResponse] {
	return func(ctx context.Context) (*hello.ListGreetingsResponse, error) {
		limit := req.Limit
		if limit <= 0 {
			limit = 10
		}

		// Query Firestore
		queryOp := hello.LiftDBOpToNetworkOp(
			s.greetings.Where("count", ">=", 0).Limit(int(limit)).GetAll(),
		)
		greetings, err := queryOp(ctx)
		if err != nil {
			return nil, err
		}

		return &hello.ListGreetingsResponse{Greetings: greetings}, nil
	}
}
