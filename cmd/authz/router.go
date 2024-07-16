package main

import "net/http"

func NewRouter() *http.ServeMux {
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
	mux.HandleFunc("/demo/document/", documentHandler)
	mux.HandleFunc("/demo/login/", loginHandler)
	mux.HandleFunc("/demo/logout/", logoutHandler)

	return mux
}
