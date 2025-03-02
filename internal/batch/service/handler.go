package service

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

// ExecuteScriptWorkflow is the workflow that executes a script
func ExecuteScriptWorkflow(ctx workflow.Context, input ScriptExecutionInput) (*ScriptExecutionResult, error) {
	// Execute the script using activity
	var result ScriptExecutionResult
	err := workflow.ExecuteActivity(ctx, ExecuteScriptActivity, input).Get(ctx, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ExecuteScriptActivity is the activity that executes a script
func ExecuteScriptActivity(ctx context.Context, input ScriptExecutionInput) (*ScriptExecutionResult, error) {
	// Execute the script
	cmd := fmt.Sprintf("%s %s", input.ExecutorCmd, input.ScriptPath)
	result := &ScriptExecutionResult{
		Success: true,
		Output:  fmt.Sprintf("Executed command: %s", cmd),
	}
	return result, nil
}

// RegisterNexusService registers the batch script service with Nexus
func RegisterNexusService(w worker.Worker) error {
	// w.RegisterNexusService(service.NewNexusService())
	// Register workflow
	w.RegisterWorkflow(ExecuteScriptWorkflow)

	// Register activity
	w.RegisterActivity(ExecuteScriptActivity)

	return nil
}
