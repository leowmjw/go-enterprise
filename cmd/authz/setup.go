package main

import (
	"app/internal/authz"
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	"log"
)

const TQ = "example-task-queue"
const orgID = "GopherLab"

func SetupSimpleWorkflow(c client.Client) {
	// Start Workflow for Org GopherLab
	// With below combos ..
	orgID := "CrabLab"
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
	// DEBUG
	//spew.Dump(docsInit)

	fmt.Println("Start Temporal Workflow ....")
	// Start the workflow
	workflowOptions := client.StartWorkflowOptions{
		ID:        orgID,
		TaskQueue: TQ,
	}
	name := "World"
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions,
		authz.SimpleWorkflow,
		authz.WFDemoInput{
			Name: name,
			Docs: docsInit,
		})
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	fmt.Print("Starting workflow for Org ", orgID)
	// Pass in the authz mdoel ..
	// If workflow already started .. no need to reinit ..

	fmt.Println("WF:", we.GetRunID())
	// Get the workflow result
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable to get workflow result", err)
	}

	fmt.Println("Workflow result:", result)
	return
}

// SetupActionWorkflow demos an action happening .. and signalling ..
func SetupActionWorkflow(c client.Client) {
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
		authz.Document{
			ID:      "secret/salary.doc",
			Owner:   "mleow",
			Content: "Lotsa Moolah!!",
		},
	}
	usersInit := []string{"bob", "mleow"}
	// DEBUG
	//spew.Dump(docsInit)

	fmt.Println("Start Temporal Workflow ==> ActionWorkflow")
	// Start the workflow - ActionWorkflow
	workflowOptions := client.StartWorkflowOptions{
		ID:        orgID,
		TaskQueue: TQ,
	}
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions,
		authz.ActionWorkflow,
		authz.WFDemoInput{
			Users: usersInit,
			Docs:  docsInit,
		})
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	fmt.Print("Started workflow for Org ", orgID, " RunID: ", we.GetRunID())
	// Below for DEBUG of flow; shoul all be triggred in the handler !!!
	// =================================>
	//// Delay to send action signal ..
	//time.Sleep(time.Second * 15)
	//saerr := c.SignalWorkflow(context.Background(), orgID, we.GetRunID(), "actionSignal", authz.Actions{
	//	GetAdminElevated: true,
	//})
	//if saerr != nil {
	//	log.Fatalln("Unable to signal workflow", saerr)
	//}
	//// Delay ,.. then terminate ...
	//time.Sleep(time.Second * 30)
	//fmt.Println("AFTER 30s!!! =====>> ****")
	//serr := c.SignalWorkflow(context.Background(), orgID, we.GetRunID(), "terminateSignal", authz.Actions{
	//	GetAdminElevated: true,
	//})
	//if serr != nil {
	//	log.Fatalln("Unable to signal workflow", serr)
	//}

	return
}
