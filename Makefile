.PHONY: api test-api research help

# Default target
help:
	@echo "Available commands:"
	@echo "  make api       - Start the API server on :8080"
	@echo "  make test-api  - Test the API with a sample request"
	@echo "  make research  - Get research report (decoded)"

# Start the API server
api:
	go run cmd/main.go

# Test the API
test-api:
	curl -s http://localhost:8080/data/bloomberg | jq .

# Get research report
research:
	@curl -s -X POST http://localhost:8080/research | jq -r '.research' | base64 -d
