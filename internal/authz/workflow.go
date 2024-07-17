package authz

import (
	"context"
	"fmt"
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

// GreetActivity .. is dummy activity ..
func GreetActivity(ctx context.Context, name string) (string, error) {
	// Simulate some work with a sleep
	time.Sleep(2 * time.Second)

	result := fmt.Sprintf("Hello, %s!", name)
	return result, nil
}
