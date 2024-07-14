package main

import (
	"fmt"
	"net/http"
)

// Handlers for accessing documents
// Parameters will mention the
// Based on the identity; extract from session cookie ..

func demoHandler(w http.ResponseWriter, r *http.Request) {
	// Show login screen or if logged in .. then documents able to be accessed ..
	// Ask for owner approval .. pending ..
	// See public doc ..
	c, err := r.Cookie("ID")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Redirect(w, r, "/demo/login/", http.StatusFound)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	//fmt.Println("ID:", c.Value)
	fmt.Fprintf(w, "Goodbye Cruel World!! See ya %s", c.Value)
	return
}
