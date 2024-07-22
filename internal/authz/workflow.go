package authz

import (
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
	TempElevated     bool
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
	// Common vars ..
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

	// For processing clean Termination ..
	var terminate bool
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
		c.Receive(ctx, nil)
		logger.Info("Received terminate signal")
		// Dump out state ..
		ad.debugState()
		// Simulate cleaning up .. and persisting data ..
		err := workflow.Sleep(ctx, time.Second) // TODO: Jitter based on random so can be between 300ms to 1500 ms
		if err != nil {
			logger.Error("Failed to sleep", "Error", err)
		}
		logger.Info("ActionWorkflow completed after persisting state")
		terminate = true // Can finally finsh ..
		return
	})

	// Process signals until a terminate signal is received
	for {
		// Wait for the next signal
		selector.Select(ctx)
		if !selector.HasPending() && terminate {
			break
		}
		// If has pending signal; should contiune processing so not lost signal ..
	}

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
	//spew.Dump(actions)
	var a *Activities
	if actions.TempElevated {
		cname := "tempaccess-mleow-secret-" + workflow.Now(ctx).String()
		// try out co-routine ..
		workflow.GoNamed(ctx, cname, func(ctx workflow.Context) {
			logger.Info("Inside co-routine:", cname)
			// Just in case; put a timeout ..
			ao := workflow.ActivityOptions{
				StartToCloseTimeout: time.Second * 10,
			}
			ctx = workflow.WithActivityOptions(ctx, ao)
			err := workflow.ExecuteActivity(ctx, a.TempAccessActivity, "mleow", "secret/secretz.doc").Get(ctx, nil)
			if err != nil {
				logger.Error("TempAccessActivity failed.", "Error", err)
				return
			}
			// Disable it after 1 min
			workflow.Sleep(ctx, time.Second*30)
			xerr := workflow.ExecuteActivity(ctx, a.RemoveAccessActivity, "mleow", "secret/secretz.doc").Get(ctx, nil)
			if xerr != nil {
				logger.Error("RemoveAccessActivity failed.", "Error", xerr)
				return
			}
		})
	}

	return
}

// ApprovalWorkflow will wait and block till ... approve or rejected ..  ID is docID ..
func ApprovalWorkflow(ctx workflow.Context, input WFDemoInput) error {
	return nil
}
