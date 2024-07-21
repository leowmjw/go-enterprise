package authz

import (
	"github.com/davecgh/go-spew/spew"
	"go.temporal.io/sdk/workflow"
	"time"
)

type WFDemoInput struct {
	Name  string
	Users []string
	Docs  []Document
}

type Actions struct {
	CheckApproval    bool
	GetAdminElevated bool
	AddPermission    bool
	RemovePermission bool
}

type WFDemoOutput struct {
	Content string
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

// ActionWorkflow - Process loop waiting ..
func ActionWorkflow(ctx workflow.Context, input WFDemoInput) error {
	// For processing clean Termination ..
	var stopSignalSent, stopSignalProcessed bool

	logger := workflow.GetLogger(ctx)
	wfInfo := workflow.GetInfo(ctx)
	workflowID := wfInfo.WorkflowExecution.ID
	logger.Info("ActionWorkflow started", "WorkflowID", workflowID)

	// Setup the first time ..
	// If have previous state; reload here ... and restart any process ..
	ad, naerr := NewAuthzDemo("", "")
	if naerr != nil {
		logger.Error("NewAuthzDemo failed.", "Error", naerr)
		return naerr
	}
	// Init data ..
	ad.users = input.Users
	ad.docs = input.Docs
	serr := ad.setupTuples()
	if serr != nil {
		logger.Error("SetupTuples failed.", "Error", serr)
		return serr
	}
	// Define signals
	var actions Actions
	signalChan := workflow.GetSignalChannel(ctx, "actionSignal")
	terminateChan := workflow.GetSignalChannel(ctx, "terminateSignal")

	// Use selector if got multiple Signal Types to handle ..
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
		stopSignalSent = true
		// Dump out state ..
		ad.debugState()
		// Simulate cleaning up .. and persisting data ..
		err := workflow.Sleep(ctx, time.Second) // TODO: Jitter based on random so can be between 300ms to 1500 ms
		if err != nil {
			logger.Error("Failed to sleep", "Error", err)
		}
		logger.Info("ActionWorkflow completed after persisting state")
		stopSignalProcessed = true // Can finally finsh ..
		return
	})

loop:
	// Wait for a terminate signal or action signal
	selector.Select(ctx)
	if !stopSignalSent {
		// Back to next signal ..
		logger.Info("Got ACTION! Back to waiting loop!!!")
		goto loop
	}
	// Block until Termination process is completed otherwise might get error like:
	//		" Workflow has unhandled signals"
	ok, err := workflow.AwaitWithTimeout(ctx, time.Minute, func() bool {
		return stopSignalProcessed == true
	})
	if err != nil {
		logger.Error("UNEXPECTED ERR:", err)
		return err
	} else if !ok {
		logger.Error("Timed out waiting for actions to be processed")
		return nil
	}
	// TODO: Process any remaining signals ... before shut it down ..
	
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
