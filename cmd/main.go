package main

import (
	"log"
	"os"

	"github.com/davidwrenner/replay-2026-hackathon/activities"
	"github.com/davidwrenner/replay-2026-hackathon/api"
	"github.com/davidwrenner/replay-2026-hackathon/pkg/workflowext"
	"github.com/davidwrenner/replay-2026-hackathon/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const TaskQueue = "research-task-queue"

func main() {
	addr := os.Getenv("API_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	temporalAddr := os.Getenv("TEMPORAL_ADDR")
	if temporalAddr == "" {
		temporalAddr = "localhost:7233"
	}

	// Create Temporal client
	c, err := client.Dial(client.Options{
		HostPort: temporalAddr,
	})
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
	defer c.Close()

	// Create worker
	w := worker.New(c, TaskQueue, worker.Options{})

	// Register workflows
	w.RegisterWorkflow(workflows.ResearchWorkflow)

	// Register activities with dependencies
	acts := &activities.Activities{
		ReportData: api.GetReportData(),
	}
	w.RegisterActivity(acts)
	w.RegisterActivity(workflowext.NoOpActivity)

	// Start worker in background
	go func() {
		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Fatalf("Failed to start worker: %v", err)
		}
	}()

	log.Printf("Temporal worker started on task queue: %s", TaskQueue)

	// Start API server
	server := api.NewServer(addr, c)
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
