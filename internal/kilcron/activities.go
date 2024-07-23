package kilcron

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/exp/rand"
	"sync"
	"time"
)

// Thread-safe counter
var (
	requestCount int
	countMutex   sync.Mutex
)

// MockHTTPCall simulates stable HTTP endpoint
func MockHTTPCall() error {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(171)+30))
	return nil
}

// MockFlakyHTTPCall simulates a flaky HTTP endpoint
func MockFlakyHTTPCall() error {
	countMutex.Lock()
	requestCount++
	currentCount := requestCount
	countMutex.Unlock()

	if currentCount <= 3 {
		fmt.Println("BEFORE:", time.Now())
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1001)+500))
		fmt.Println("AFTER:", time.Now())
		// 90% error
		if rand.Float32() < 0.9 {
			return errors.New("internal server error")
		}
		// OK but slow ..
	} else {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(471)+30))
		// 90% is OK
		if rand.Float32() >= 0.9 {
			return errors.New("internal server error")
			return nil
		}
	}
	// All OK ..
	return nil
}

func MakePayment(ctx context.Context, paymentID string) error {
	fmt.Println("MakePayment: ID:", paymentID)
	// Below to test flaky calls ..
	return MockFlakyHTTPCall()
	//return MockHTTPCall()
	return nil
}
