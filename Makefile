.PHONY: generate build run-server run-client run-web clean install-deps all dev

# Generate protobuf code and copy to web
generate:
	buf generate
	rm -rf web/src/gen
	mkdir -p web/src/gen
	cp -r gen/ts/* web/src/gen/

# Build binaries
build:
	go build -o bin/server ./cmd/server
	go build -o bin/client ./cmd/client

# Run server
run-server:
	GOOGLE_CLOUD_PROJECT=graphql-category-db go run ./cmd/server

# Run Go client
run-client:
	go run ./cmd/client

# Run web frontend
run-web:
	cd web && npx vite

# Run both server and web (requires tmux or run in separate terminals)
dev:
	@echo "Run in separate terminals:"
	@echo "  Terminal 1: make run-server"
	@echo "  Terminal 2: make run-web"

# Install dependencies
install-deps:
	go mod tidy
	cd web && npm install

# Clean generated and built files
clean:
	rm -rf bin/
	rm -rf gen/
	rm -rf web/src/gen/
	rm -rf web/dist/

# Full rebuild
all: clean generate build
	@echo "Build complete!"
