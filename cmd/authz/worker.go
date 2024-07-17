package main

import (
	"app/internal/authz"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// NewDemoWorker ...
func NewDemoWorker(c client.Client) worker.Worker {
	// Create the Temporal client
	//c, err := client.NewLazyClient(client.Options{})
	//if err != nil {
	//	log.Fatalln("Unable to create Temporal client", err)
	//}
	//defer c.Close()

	// Create a worker that listens on the task queue and hosts the workflow and activity functions
	w := worker.New(c, TQ, worker.Options{})

	w.RegisterWorkflow(authz.SimpleWorkflow)
	w.RegisterActivity(authz.GreetActivity)

	return w
}
