package authz

import (
	"context"
	"fmt"
	"time"
)

type Activities struct {
	As AuthStore
}

// GreetActivity .. is dummy activity ..
func GreetActivity(ctx context.Context, name string) (string, error) {
	// Simulate some work with a sleep
	time.Sleep(2 * time.Second)

	result := fmt.Sprintf("Hello, %s!", name)
	return result, nil
}

func (a *Activities) TempAccessActivity(ctx context.Context, user, document string) error {
	fmt.Println("Inside TempAccessActivity .. let;s see if can see OpenFGA ..")
	a.As.InitDemo("")
	return nil
}
