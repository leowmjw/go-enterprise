package main

import (
	"context"
	"fmt"
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
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Define a handler function
	h1 := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	}

	// Attach handler function to the ServeMux
	mux.HandleFunc("/", h1)

	// Create the Server using the new ServeMux
	server := &http.Server{
		Addr:    ":8888",
		Handler: mux,
	}

	// Prepare for handling signals
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Running the HTTP server in a go routine
	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println("Server error:", err)
		}
	}()

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
