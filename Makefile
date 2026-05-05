.PHONY: api help

# Default target
help:
	@echo "Available commands:"
	@echo "  make api    - Start the API server on :8080"

# Start the API server
api:
	go run cmd/main.go
