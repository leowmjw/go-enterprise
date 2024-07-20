package authz

import (
	"github.com/davecgh/go-spew/spew"
	"go.temporal.io/sdk/workflow"
	"time"
)

type WFDemoInput struct {
	Name string
	Docs []Document
}

// SimpleWorkflow is dummy workflow ...
func SimpleWorkflow(ctx workflow.Context, input WFDemoInput) (string, error) {
	// Workflow ..
	// DEBUG
	//spew.Dump(input.Docs)
	// Start with the Authorization Models initailized
	// Block until get signal to complete ...

	// Workflow code goes here
	logger := workflow.GetLogger(ctx)
	logger.Info("SimpleWorkflow started", "name", input.Name)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 10,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	err := workflow.ExecuteActivity(ctx, GreetActivity, input.Name).Get(ctx, &result)
	if err != nil {
		logger.Error("GreetActivity failed.", "Error", err)
		return "", err
	}

	logger.Info("SimpleWorkflow completed", "result", result)
	return result, nil
}

type Actions struct {
	CheckApproval    bool
	GetAdminElevated bool
	AddPermission    bool
	RemovePermission bool
}

// ActionWorkflow - Process loop waiting ..
func ActionWorkflow(ctx workflow.Context, input WFDemoInput) error {
	logger := workflow.GetLogger(ctx)
	wfInfo := workflow.GetInfo(ctx)
	workflowID := wfInfo.WorkflowExecution.ID
	logger.Info("ActionWorkflow started", "WorkflowID", workflowID)

	// Setup the first time ..
	// If have previous state; reload here ... and restart any process ..

	// Define signals
	var actions Actions
	signalChan := workflow.GetSignalChannel(ctx, "actionSignal")
	terminateChan := workflow.GetSignalChannel(ctx, "terminateSignal")

	selector := workflow.NewSelector(ctx)

	// Handling Actions ..
	selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &actions)
		logger.Info("Received signal", "actions", actions)
		handleActions(ctx, actions)
	})

	// Handling Termination + state saving mechanism ..
	selector.AddReceive(terminateChan, func(c workflow.ReceiveChannel, more bool) {
		logger.Info("Received terminate signal")
		// Simulate cleaning up .. and persisting data ..
		err := workflow.Sleep(ctx, time.Second) // TODO: Jitter based on random so can be between 300ms to 1500 ms
		if err != nil {
			logger.Error("Failed to sleep", "Error", err)
		}
		logger.Info("ActionWorkflow completed after persisting state")
		return
	})

	// Wait for a terminate signal or action signal
	selector.Select(ctx)

	logger.Info("ActionWorkflow completed")
	return nil
}

func handleActions(ctx workflow.Context, actions Actions) {
	// Implement action handling logic here
	logger := workflow.GetLogger(ctx)
	wfInfo := workflow.GetInfo(ctx)
	if wfInfo.GetCurrentHistoryLength() > 100 {
		logger.Info("History length is too large.")
		// Below only for when getting close to limit run ..?
		//err := workflow.NewContinueAsNewError(ctx, ActionWorkflow)
		//if err != nil {
		//	logger.Error(err.Error())
		//	return err
		//}
	}
	wfEx := wfInfo.WorkflowExecution
	// DEBUIG
	logger.Info("Inside handleActions:",
		"ID", wfEx.ID, "RunID", wfEx.RunID,
		"actions", actions,
	)
	spew.Dump(actions)
}

// ApprovalWorkflow will wait and block till ... approve or rejected ..  ID is docID ..
func ApprovalWorkflow(ctx workflow.Context, input WFDemoInput) error {
	return nil
}
