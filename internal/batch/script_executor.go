package batch

import (
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"app/internal/batch/service"
)

// List of random sites to call
var randomSites = []string{
	"https://example.com",
	"https://github.com",
	"https://google.com",
	"https://wikipedia.org",
	"https://stackoverflow.com",
	"https://go.dev",
	"https://reddit.com",
	"https://twitter.com",
	"https://amazon.com",
	"https://microsoft.com",
}

// ExecuteScript runs a script using the specified executor
func ExecuteScript(ctx context.Context, input service.ScriptExecutionInput) (*service.ScriptExecutionResult, error) {
	// Check if this is a special HTTP request batch
	if input.ExecutorCmd == "uv" && input.APIFunction == "http_batch" {
		return executeHttpBatch(ctx)
	}

	// Regular script execution
	scriptDir := filepath.Join("internal/batch/scripts", filepath.Clean(input.ScriptPath))
	
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

// executeHttpBatch executes a batch of HTTP requests using the uv command
func executeHttpBatch(ctx context.Context) (*service.ScriptExecutionResult, error) {
	// Initialize random generator with seed
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// String builder for results
	var results strings.Builder
	results.WriteString("HTTP Batch Results:\n\n")
	
	// Execute 10 HTTP requests to random sites
	success := true
	for i := 0; i < 10; i++ {
		// Select a random site
		site := randomSites[r.Intn(len(randomSites))]
		
		// Create the uv command with the site URL
		cmd := exec.CommandContext(ctx, "uv", "-X", "GET", site)
		
		// Execute the command
		output, err := cmd.CombinedOutput()
		
		// Record the results
		results.WriteString(fmt.Sprintf("Request %d: %s\n", i+1, site))
		if err != nil {
			success = false
			results.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
		} else {
			// Truncate the output to avoid extremely long responses
			outputStr := string(output)
			if len(outputStr) > 1000 {
				outputStr = outputStr[:1000] + "... (truncated)"
			}
			results.WriteString(fmt.Sprintf("Status: Success\nResponse snippet: %s\n", outputStr))
		}
		results.WriteString("\n---\n\n")
		
		// Add a small delay between requests
		time.Sleep(100 * time.Millisecond)
	}
	
	return &service.ScriptExecutionResult{
		Success: success,
		Output:  results.String(),
	}, nil
}
