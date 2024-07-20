package main

import (
	"context"
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

	// Document Lists

	// Pending Approvers ..
	return
}

func documentHandler(w http.ResponseWriter, r *http.Request) {
	// Get our Demo Workflow ..
	wfr := c.GetWorkflow(context.Background(), orgID, "")
	wfr.GetRunID()

	// WorkflowID: <username>-approver
	// WorkflowID: <docID>

	// Check if action is happening ... after done redirect back ..
	q := r.URL.Query()
	if q.Has("action") {

		switch q.Get("action") {
		case "approve":
			// Redirect back it itself?

		case "reject":
			// Redirect back to top level??

		case "view":
		// If no document .. BadRequest
		// Check if got viewer access or not ..
		// if yes, show secrets .. else naughty! can for access

		case "kil":
			err := c.SignalWorkflow(context.Background(), orgID, "", "terminateSignal", true)
			if err != nil {
				fmt.Println("ERR: ", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	// Show Workflow current status under care??
	fmt.Fprintf(w, "Nothing to see here .. docs")
	return
}
