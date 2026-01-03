.PHONY: generate build run-server run-client clean

# Generate protobuf code
generate:
	buf generate

# Build binaries
build:
	go build -o bin/server ./cmd/server
	go build -o bin/client ./cmd/client

# Run server
run-server:
	go run ./cmd/server

# Run client
run-client:
	go run ./cmd/client

# Clean generated and built files
clean:
	rm -rf bin/
	rm -rf gen/go/hello/*.pb.go
	rm -rf gen/go/hello/*_category.go
	rm -rf gen/go/hello/helloconnect/
	rm -rf gen/go/category/
