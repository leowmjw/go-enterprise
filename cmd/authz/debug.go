package main

import (
	"context"
	"fmt"
	"net/http"
)

func renderDefault() string {
	result := "<h3><strong>ACCESS MATRIX</strong></h3>"
	result += "<div>"
	// Get current model
	// Test access for all the users .. print out the report ..
	// Grant .. and check ...
	// Print link here ...
	users := []string{"bob", "mleow"}
	docs := []string{"public/welcome.doc", "secret/secretz.doc"}
	for _, user := range users {
		for _, doc := range docs {
			result += doc
			ok, _ := as.CanViewDocument(user, doc)
			if ok {
				result += " - YES "
			} else {
				result += " - NO "
			}
			result += "<br/>"
		}
	}
	result += "</div>"
	result += `
<div>
	<a href=\"/demo/debug/?action=temp\">Grant Temp Access</a>
	<a href=\"/demo/debug/?action=kil\">Terminate</a>
</div>`

	return result
}

func debugAccessHandler(w http.ResponseWriter, r *http.Request) {
	// Get our Demo Workflow ..
	wfr := c.GetWorkflow(context.Background(), orgID, "")
	wfr.GetRunID()

	// WorkflowID: <username>-approver
	// WorkflowID: <docID>

	// Check if action is happening ... after done redirect back ..
	q := r.URL.Query()
	if q.Has("action") {
		switch q.Get("action") {
		case "temp":
			// Temp access for 2 mins??

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
	result := renderDefault()
	// Remove it .. check in 30s

	fmt.Fprintf(w, result)
	return
}
