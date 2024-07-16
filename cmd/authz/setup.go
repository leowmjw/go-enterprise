package main

import (
	"app/internal/authz"
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

// Setup Temporal Client??

func init() {

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
	spew.Dump(docsInit)
	fmt.Print("Starting workflow for Org ", orgID)
	// Pass in the authz mdoel ..
	// If workflow already started .. no need to reinit ..
}
