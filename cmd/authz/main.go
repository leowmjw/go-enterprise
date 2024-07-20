package main

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go.temporal.io/sdk/client"
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

var c client.Client

func Run() {
	// HTTP Server Setup ..
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
	// END HTTP Server ================>

	// Create the Temporal client
	var err error
	c, err = client.NewLazyClient(client.Options{})
	if err != nil {
		spew.Dump(err)
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Setup the Sanity Test Scenario ..
	go SetupSimpleWorkflow(c)
	// Actual Demo Scenario ..
	//go SetupActionWorkflow(c)

	// Running the Temporal Worker in a go routine ..
	go SetupTemporalWorker(c)

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
}
