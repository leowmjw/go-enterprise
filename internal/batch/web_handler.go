package batch

import (
	"embed"
	"html/template"
	"net/http"

	"app/internal/batch/service"
	"go.temporal.io/sdk/client"
)

//go:embed templates/*
var templateFS embed.FS

// WebHandler handles HTTP requests for script execution
type WebHandler struct {
	client    client.Client
	templates *template.Template
}

// NewWebHandler creates a new web handler
func NewWebHandler(c client.Client) (*WebHandler, error) {
	tmpl, err := template.ParseFS(templateFS, "templates/*.html")
	if err != nil {
		return nil, err
	}

	return &WebHandler{
		client:    c,
		templates: tmpl,
	}, nil
}

// RegisterRoutes registers the web handler routes
func (h *WebHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/execute", h.handleExecute)
	mux.HandleFunc("/execute/submit", h.handleSubmit)
	mux.HandleFunc("/execute/status/", h.handleStatus)
}

func (h *WebHandler) handleExecute(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "execute.html", nil)
}

func (h *WebHandler) handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	input := service.ScriptExecutionInput{
		APIFunction: r.FormValue("apiFunction"),
		ScriptPath:  r.FormValue("scriptPath"),
		ExecutorCmd: r.FormValue("executorCmd"),
		NexusPath:   r.FormValue("nexusPath"),
	}

	// Start workflow
	workflowOptions := client.StartWorkflowOptions{
		ID:        "script-execution-" + input.ScriptPath,
		TaskQueue: "batch-demo",
	}

	we, err := h.client.ExecuteWorkflow(r.Context(), workflowOptions, service.ExecuteScriptWorkflow, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return partial update for htmx
	h.templates.ExecuteTemplate(w, "status_partial.html", map[string]interface{}{
		"WorkflowID": we.GetID(),
		"Status":     "Started",
	})
}

func (h *WebHandler) handleStatus(w http.ResponseWriter, r *http.Request) {
	workflowID := r.URL.Path[len("/execute/status/"):]
	
	// Get workflow execution
	we := h.client.GetWorkflow(r.Context(), workflowID, "")

	// Get workflow result
	var result service.ScriptExecutionResult
	statusText := "Running"

	err := we.Get(r.Context(), &result)
	if err == nil {
		if result.Success {
			statusText = "Completed"
		} else {
			statusText = "Failed"
		}
	}

	// Return partial update for htmx
	h.templates.ExecuteTemplate(w, "status_partial.html", map[string]interface{}{
		"WorkflowID": workflowID,
		"Status":     statusText,
		"Result":     &result,
	})
}
