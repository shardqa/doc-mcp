.PHONY: test build clean run

build:
	@echo "Building doc-mcp..."
	go build -o doc-mcp .
	@echo "✅ Build successful! Binary: ./doc-mcp"

run: build
	./doc-mcp

clean:
	@echo "Cleaning build artifacts..."
	rm -f doc-mcp coverage.out
	@echo "✅ Clean complete"

test:
	@echo "Running all Go tests..."
	go test ./... 

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out 