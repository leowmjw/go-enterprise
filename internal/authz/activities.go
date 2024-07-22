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
	//a.As.AllowView

	// If already got access .. to handle it??
	err := a.As.AddViewRelationship(user, document)
	if err != nil {
		fmt.Println("Error adding view relationship. ERR:", err)
		// Error will cause it to retry .. how long?
		// Should catch exist error .. and move on ..
		//return err
	}
	return nil
}

func (a *Activities) RemoveAccessActivity(ctx context.Context, user, document string) error {
	fmt.Println("Inside RemoveAccessActivity .. clean up ..")
	//a.As.AllowView
	err := a.As.RemoveViewRelationship(user, document)
	if err != nil {
		fmt.Println("Error removing view relationship. ERR:", err)
	}
	return nil
}
