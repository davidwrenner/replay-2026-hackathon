.PHONY: api test-api help

# Default target
help:
	@echo "Available commands:"
	@echo "  make api       - Start the API server on :8080"
	@echo "  make test-api  - Test the API with a sample request"

# Start the API server
api:
	go run cmd/main.go

# Test the API
test-api:
	curl -s http://localhost:8080/data/bloomberg | jq .
