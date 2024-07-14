package main

import (
	"app/internal/authz"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var as authz.AuthStore

func init() {
	apiURL := os.Getenv("FGA_API_URL")
	as = authz.NewAuthStore(apiURL)
	err := as.InitDemo("")
	if err != nil {
		fmt.Println("Error initializing auth store. ERR:", err.Error())
	} else {
		fmt.Println("SUCCESS!! Init Authz Server!!")
	}
}

func main() {
	fmt.Println("Demo Server ...")
	Run()
}

func Run() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Define a handler function
	defaultHandler := func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "Hello World")
		http.NotFound(w, r)
		return
	}

	// Attach handler function to the ServeMux
	mux.HandleFunc("/", defaultHandler)
	mux.HandleFunc("/demo/", demoHandler)
	mux.HandleFunc("/demo/login/", loginHandler)
	mux.HandleFunc("/demo/logout/", logoutHandler)

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

	// Shutdown Temporal Worker ...
}
