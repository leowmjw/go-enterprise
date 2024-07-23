package main

import (
	"app/internal/kilcron"
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	"net/http"
)

const TQ = "kilcron-task-queue"
const orgID = "GopherPayNET"

// debugAccessHandler will start flaky cron ..
func debugAccessHandler(w http.ResponseWriter, r *http.Request) {
	payID := "GoBux"
	wfr, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
		ID:        orgID,
		TaskQueue: TQ,
	}, kilcron.PaymentWorkflow, payID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to start workflow: %v", err)
		return
	}
	render := "WF Started - ID: " + wfr.GetID() + " RunID: " + wfr.GetRunID()
	// DEBUG
	fmt.Println(render)
	fmt.Fprintf(w, render)
	return
}
