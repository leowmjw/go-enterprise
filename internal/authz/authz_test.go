package authz

import (
	"os"
	"testing"
)

func TestCheckPermission(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{"Case #1"},
	}
	// Setup ...
	os.Setenv("FGA_API_URL", "http://localhost:8080")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckPermission()
		})
	}
}
