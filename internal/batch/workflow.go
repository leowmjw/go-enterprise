package batch

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// BatchWorkflow represents the main workflow that runs multiple scenarios
func BatchWorkflow(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 2,
			InitialInterval: time.Second,
			BackoffCoefficient: 2.0,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Run 10 instances of each scenario
	for i := 0; i < 10; i++ {
		// Run scenarios in parallel
		future1a := workflow.ExecuteActivity(ctx, Scenario1a)
		future1b := workflow.ExecuteActivity(ctx, Scenario1b)
		future2a := workflow.ExecuteActivity(ctx, Scenario2a)
		future2b := workflow.ExecuteActivity(ctx, Scenario2b)

		// Wait for all scenarios to complete
		if err := future1a.Get(ctx, nil); err != nil {
			workflow.GetLogger(ctx).Error("Scenario1a failed", "error", err)
		}
		if err := future1b.Get(ctx, nil); err != nil {
			workflow.GetLogger(ctx).Error("Scenario1b failed", "error", err)
		}
		if err := future2a.Get(ctx, nil); err != nil {
			workflow.GetLogger(ctx).Error("Scenario2a failed", "error", err)
		}
		if err := future2b.Get(ctx, nil); err != nil {
			workflow.GetLogger(ctx).Error("Scenario2b failed", "error", err)
		}
	}

	return nil
}
