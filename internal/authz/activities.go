package authz

import (
	"context"
	"fmt"
	"time"
)

// GreetActivity .. is dummy activity ..
func GreetActivity(ctx context.Context, name string) (string, error) {
	// Simulate some work with a sleep
	time.Sleep(2 * time.Second)

	result := fmt.Sprintf("Hello, %s!", name)
	return result, nil
}
