package service

const (
    BatchServiceName = "batch-script-service"
    ExecuteScriptOperationName = "execute-script"
)

// ScriptExecutionInput represents the input for script execution
type ScriptExecutionInput struct {
    APIFunction   string `json:"apiFunction"`   
    ScriptPath    string `json:"scriptPath"`    
    ExecutorCmd   string `json:"executorCmd"`   
    NexusPath     string `json:"nexusPath"`     
}

// ScriptExecutionResult represents the result of script execution
type ScriptExecutionResult struct {
    Success      bool   `json:"success"`
    Output       string `json:"output"`
    ErrorMessage string `json:"errorMessage,omitempty"`
}
