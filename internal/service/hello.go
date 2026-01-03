package service

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/example/hello-ct/gen/go/hello"
)

type helloService struct {
	greetings sync.Map
	counter   atomic.Int64
}

// NewHelloService creates a HelloServiceK using CT patterns
func NewHelloService() *hello.HelloServiceK {
	svc := &helloService{}

	return &hello.HelloServiceK{
		SayHello:      svc.sayHello,
		ListGreetings: svc.listGreetings,
	}
}

func (s *helloService) sayHello(req *hello.SayHelloRequest) hello.NetworkOp[*hello.SayHelloResponse] {
	return func(ctx context.Context) (*hello.SayHelloResponse, error) {
		count := s.counter.Add(1)

		greeting := &hello.Greeting{
			Id:      fmt.Sprintf("greeting-%d", count),
			Message: fmt.Sprintf("Hello, %s!", req.Name),
			Count:   int32(count),
		}

		s.greetings.Store(greeting.Id, greeting)

		return &hello.SayHelloResponse{
			Greeting: greeting,
		}, nil
	}
}

func (s *helloService) listGreetings(req *hello.ListGreetingsRequest) hello.NetworkOp[*hello.ListGreetingsResponse] {
	return func(ctx context.Context) (*hello.ListGreetingsResponse, error) {
		var greetings []*hello.Greeting

		limit := req.Limit
		if limit <= 0 {
			limit = 10
		}

		s.greetings.Range(func(key, value any) bool {
			if len(greetings) >= int(limit) {
				return false
			}
			greetings = append(greetings, value.(*hello.Greeting))
			return true
		})

		return &hello.ListGreetingsResponse{
			Greetings: greetings,
		}, nil
	}
}
