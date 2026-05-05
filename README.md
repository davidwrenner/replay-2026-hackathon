# replay-2026 hackathon

Kalshi Bot — AI-powered market research and trading strategy tool built with React, Go, and Temporal workflows.

## Architecture

- **Frontend** — React + TypeScript (Vite), renders MDX reports with Vega-Lite charts
- **Backend** — Go API server on `:8080`, orchestrates research via Temporal workflows
- **Workflows** — Temporal workflows coordinate data source queries and report generation

## Local Setup

### Prerequisites

- [Go](https://go.dev/dl/) 1.21+
- [Node.js](https://nodejs.org/) 18+
- [Temporal CLI](https://docs.temporal.io/cli)

### 1. Start Temporal

```shell
temporal server start-dev
# Temporal Server:  localhost:7233
# Temporal UI:      http://localhost:8233
```

### 2. Start the API Server

```shell
make api
# API server on :8080
```

### 3. Start the Frontend

```shell
make ui
# Frontend on http://localhost:5173 (proxies /research → :8080)
```

## Available Commands

| Command | Description |
|---|---|
| `make help` | Show all available commands |
| `make api` | Start the Go API server on `:8080` |
| `make ui` | Start the React dev server on `:5173` |
| `make test-api` | Test the API with a sample request |
| `make research` | Get a decoded research report via curl |
