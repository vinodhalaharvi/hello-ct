# Hello CT

Minimal starter template using Category Theory protoc plugin with Connect-go.

## Prerequisites

```bash
# Install buf
brew install bufbuild/buf/buf

# Install protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install connectrpc.com/connect/cmd/protoc-gen-connect-go@latest

# Install CT plugin from local source
cd ../protoc-gen-category
go install ./cmd/protoc-gen-category

# Verify all are in PATH
which protoc-gen-go protoc-gen-connect-go protoc-gen-category
```

## Directory Layout

```
hello-ct/
├── proto/
│   ├── buf.yaml
│   ├── category/options.proto    # CT plugin options (from plugin)
│   └── hello/hello.proto         # Your service
├── gen/go/                       # Generated (after buf generate)
├── internal/service/hello.go     # Service implementation
├── cmd/server/main.go
├── cmd/client/main.go
├── buf.work.yaml
├── buf.gen.yaml
└── go.mod                        # Has replace directive for local plugin
```

## Setup

```bash
# 1. Generate code
buf generate

# 2. Tidy modules
go mod tidy

# 3. Build
go build ./...
```

## Run

```bash
# Terminal 1
go run ./cmd/server

# Terminal 2
go run ./cmd/client
```

## What the CT Plugin Generates

### Effect Types
```go
type NetworkOp[A any] func(context.Context) (A, error)
type DBOp[A any] func(context.Context) (A, error)
```

### Monad Operations
```go
PureNetwork(value)                    // Lift value
FMapNetwork(f, op)                    // Transform result
BindNetwork(op, f)                    // Chain operations
TraverseNetworkParallel(items, f)    // Parallel map
```

### Kleisli Service
```go
type HelloServiceK struct {
    SayHello      func(*SayHelloRequest) NetworkOp[*SayHelloResponse]
    ListGreetings func(*ListGreetingsRequest) NetworkOp[*ListGreetingsResponse]
}
```

### Middleware
```go
ApplyHelloServiceMiddleware(svc, &HelloServiceMiddleware{
    SayHello: myMiddleware,
})
```

### Monoid
```go
GreetingMonoid.Empty()
GreetingMonoid.Combine(g1, g2)
FoldMap(greetings, Id[*Greeting](), GreetingMonoid)
```

### Connect Bridge
```go
// Client: Connect -> Kleisli
client := NewHelloServiceKFromConnect(connectClient)

// Server: Kleisli -> Connect
handler := NewHelloServiceConnectHandler(kleisliService)
```

### Resilience
```go
RetryNetwork(op, 3, 100*time.Millisecond, 2*time.Second)
FallbackNetwork(primary, fallback)
TimeoutNetwork(op, 5*time.Second)
```

## Customizing

### Add Firestore
```protobuf
message User {
    option (category.category) = {
        firestore_bridge: true
    };
}
```

### Add Stripe
```protobuf
option (category.category_file) = {
    stripe: { plans: ["free", "pro"] }
};

message User {
    option (category.category) = {
        stripe_customer: true
    };
}
```

### Add Multi-tenancy
```go
ctx = hello.WithTenantID(ctx, "tenant-123")
// All Firestore ops auto-scope to tenant
```
# hello-ct
