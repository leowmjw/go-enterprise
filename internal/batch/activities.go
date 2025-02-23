package batch

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// FatalError represents an error that should not be retried
type FatalError struct {
	msg string
}

func (e *FatalError) Error() string {
	return e.msg
}

// Scenario1a attempts to open a non-existent file (80% failure rate)
func Scenario1a(ctx context.Context) error {
	if rand.Float64() < 0.8 {
		_, err := os.Open("/nonexistent/file")
		return &FatalError{msg: fmt.Sprintf("Fatal error opening file: %v", err)}
	}
	return nil
}

// Scenario1b mocks an endpoint with failures and timeouts (20% failure rate)
func Scenario1b(ctx context.Context) error {
	if rand.Float64() < 0.2 {
		// Simulate timeout or failure
		if rand.Float64() < 0.5 {
			time.Sleep(3 * time.Second) // Exceeds 2s timeout
			return errors.New("timeout error")
		}
		return errors.New("endpoint failure")
	}
	return nil
}

// Scenario2a attempts to connect to different ports (80% failure rate)
func Scenario2a(ctx context.Context) error {
	if rand.Float64() < 0.8 {
		return fmt.Errorf("failed to connect to port 123456")
	}
	return nil // Simulates successful connection to port 8080
}

// Scenario2b simulates authentication failures (20% failure rate)
func Scenario2b(ctx context.Context) error {
	if rand.Float64() < 0.2 {
		return errors.New("401 Unauthorized")
	}
	return nil
}
