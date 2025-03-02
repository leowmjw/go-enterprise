package batch

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"app/internal/batch/service"
)



// ExecuteScript runs a script using the specified executor
func ExecuteScript(ctx context.Context, input service.ScriptExecutionInput) (*service.ScriptExecutionResult, error) {
	// Construct full script path
	scriptDir := filepath.Join("scripts", filepath.Clean(input.ScriptPath))
	
	// Prepare command with API function as environment variable
	cmd := exec.CommandContext(ctx, input.ExecutorCmd, scriptDir)
	cmd.Env = append(cmd.Env, fmt.Sprintf("API_FUNCTION=%s", input.APIFunction))
	
	// Execute command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &service.ScriptExecutionResult{
			Success:      false,
			Output:       string(output),
			ErrorMessage: err.Error(),
		}, nil
	}
	
	return &service.ScriptExecutionResult{
		Success: true,
		Output:  string(output),
	}, nil
}
