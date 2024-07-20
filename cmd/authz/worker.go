package main

import (
	"app/internal/authz"
	"fmt"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"os"
	"os/signal"
	"syscall"
)

func SetupTemporalWorker(c client.Client) {
	fmt.Println("Run Temporal Worker ....")
	// Create a worker that listens on the task queue and hosts the workflow and activity functions
	w := worker.New(c, TQ, worker.Options{})

	w.RegisterWorkflow(authz.SimpleWorkflow)
	w.RegisterActivity(authz.GreetActivity)
	err := w.Start()
	if err != nil {
		fmt.Println("Worker error:", err)
	}

	// Prepare for handling signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	interruptSignal := <-interrupt
	fmt.Printf("TEMPORAL-WORKER: Received %s, shutting down.\n", interruptSignal)

	// Shutdown Temporal Worker ...
	fmt.Println("Stopping Temporal Worker...")
	w.Stop()
}
