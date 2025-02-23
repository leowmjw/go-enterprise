package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/internal/batch"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	fmt.Println("Temporal Batch demo ..")
	Run()
}

func Run() {
	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create Temporal client
	c, err := client.NewLazyClient(client.Options{})
	if err != nil {
		log.Fatalf("Failed to create Temporal client: %v", err)
	}
	defer c.Close()

	// Create worker variable to track worker instance
	var w worker.Worker

	// Start worker in a goroutine
	go func() {
		w = worker.New(c, "batch-demo", worker.Options{})
		w.RegisterWorkflow(batch.BatchWorkflow)
		w.RegisterActivity(batch.Scenario1a)
		w.RegisterActivity(batch.Scenario1b)
		w.RegisterActivity(batch.Scenario2a)
		w.RegisterActivity(batch.Scenario2b)

		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Printf("Worker stopped: %v", err)
		}
	}()

	// Start workflow in a goroutine
	go func() {
		workflowOptions := client.StartWorkflowOptions{
			ID:        "batch-workflow",
			TaskQueue: "batch-demo",
		}

		we, err := c.ExecuteWorkflow(ctx, workflowOptions, batch.BatchWorkflow)
		if err != nil {
			log.Printf("Failed to start workflow: %v", err)
			return
		}

		var result any
		if err := we.Get(ctx, &result); err != nil {
			log.Printf("Workflow completed with error: %v", err)
		} else {
			log.Printf("Workflow completed successfully")
		}
	}()

	// Create HTTP server
	srv := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Batch demo server running")
		}),
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on :8080")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("Shutdown signal received")

	// Create a timeout context for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Cancel the main context to stop workflow
	cancel()

	// Gracefully shutdown the HTTP server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Stop the worker
	if w != nil {
		w.Stop()
	}

	log.Println("Shutdown complete")
}
