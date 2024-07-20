package authz

import (
	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/testsuite"
	"testing"
	"time"
)

func TestActionWorkflow(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}
	env := s.NewTestWorkflowEnvironment()

	// Mock signals
	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow("actionSignal", Actions{CheckApproval: true})
	}, time.Minute)

	env.RegisterDelayedCallback(func() {
		env.SignalWorkflow("terminateSignal", true)
	}, 2*time.Minute)

	env.ExecuteWorkflow(ActionWorkflow, WFDemoInput{
		Name: "",
		Docs: nil,
	})

	// Assertions to ensure workflow handled signals and completed correctly
	assert.True(t, env.IsWorkflowCompleted())
	assert.NoError(t, env.GetWorkflowError())
}
