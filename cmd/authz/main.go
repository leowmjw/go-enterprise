package main

import (
	"app/internal/authz"
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("Demo Server ...")
	Run()
}

func Run() {
	// AuthzDemo ..
	//apiURL := os.Getenv("FGA_API_URL")
	//demo := authz.NewAuthzDemo(apiURL, "")
	// Init Temporal + OpenFGA Client
	// Start the workflow or continue ..

	// Create the Server using the new ServeMux
	server := &http.Server{
		Addr:    ":8888",
		Handler: NewRouter(),
	}

	// Running the HTTP server in a go routine
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("Server error:", err)
		}
	}()

	// Start Workflow for Org GopherLab
	// With below combos ..
	orgID := "GopherLab"
	docsInit := []authz.Document{
		authz.Document{
			ID:      "public/welcome.doc",
			Owner:   "",
			Content: "All Open!",
		},
		authz.Document{
			ID:      "secret/secretz.doc",
			Owner:   "bob",
			Content: "Secretz",
		},
	}
	// DEBUG
	//spew.Dump(docsInit)

	// Create the Temporal client
	c, err := client.NewLazyClient(client.Options{})
	if err != nil {
		spew.Dump(err)
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	go func() {
		fmt.Println("Start Temporal Workflow ....")
		// Start the workflow
		workflowOptions := client.StartWorkflowOptions{
			ID:        orgID,
			TaskQueue: TQ,
		}
		name := "World"
		we, err := c.ExecuteWorkflow(context.Background(), workflowOptions,
			authz.SimpleWorkflow,
			authz.WFDemoInput{
				name,
				docsInit,
			})
		if err != nil {
			log.Fatalln("Unable to execute workflow", err)
		}
		fmt.Print("Starting workflow for Org ", orgID)
		// Pass in the authz mdoel ..
		// If workflow already started .. no need to reinit ..

		fmt.Println("WF:", we.GetRunID())
		// Get the workflow result
		var result string
		err = we.Get(context.Background(), &result)
		if err != nil {
			log.Fatalln("Unable to get workflow result", err)
		}

		fmt.Println("Workflow result:", result)
	}()

	// Running the Temporal Worker in a go routine ..
	// passing in the clients ..
	var w worker.Worker
	go func() {
		fmt.Println("Run Temporal Worker ....")
		w = NewDemoWorker(c)
		err := w.Start()
		if err != nil {
			fmt.Println("Worker error:", err)
		}
	}()

	// Prepare for handling signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	interruptSignal := <-interrupt
	fmt.Printf("Received %s, shutting down.\n", interruptSignal)

	// Shutdown the server gracefully
	if err := server.Shutdown(context.Background()); err != nil {
		fmt.Println("Error shutting down:", err)
	} else {
		fmt.Println("Server shutdown gracefully.")
	}

	// Shutdown Temporal Worker ...
	fmt.Println("Stopping Temporal Worker...")
	w.Stop()
}
