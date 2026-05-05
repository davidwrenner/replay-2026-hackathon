.PHONY: api test-api research ui ui-install help

# Default target
help:
	@echo "Available commands:"
	@echo "  make api        - Start the API server on :8080"
	@echo "  make ui         - Install (if needed) and start the frontend dev server on :5173"
	@echo "  make test-api   - Test the API with a sample request"
	@echo "  make research   - Get research report (decoded)"

# Start the API server
api:
	go run cmd/main.go

# Start the frontend dev server (installs dependencies if needed, proxies /research to :8080)
ui:
	cd ui && if [ ! -d "node_modules" ]; then npm install; fi && npm run dev

# Test the API
test-api:
	curl -s http://localhost:8080/data/bloomberg | jq .

# Get research report
research:
	@curl -s -X POST http://localhost:8080/research | jq -r '.research' | base64 -d
