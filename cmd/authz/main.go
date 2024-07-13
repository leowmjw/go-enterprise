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

	// Create the Server
	server := &http.Server{
		Addr: ":8888",
	}

	h1 := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")

	}
	http.HandleFunc("/", h1)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Running the HTTP server
	go func() {
		server.ListenAndServe()
	}()

	interruptSignal := <-interrupt
	switch interruptSignal {
	case os.Kill:
		fmt.Println("Got SIGKILL...")
	case os.Interrupt:
		fmt.Println("Got SIGINT...")
	case syscall.SIGTERM:
		fmt.Println("Got SIGTERM...")
	}
	server.Shutdown(context.Background())
}
