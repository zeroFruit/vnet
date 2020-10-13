package test

import (
	"fmt"
	"sync"
	"time"
)

func Fatalf(format string, a ...interface{}) {
	panic(fmt.Sprintf(format, a...))
}

func ShouldSuccess(fn func() error) {
	if err := fn(); err != nil {
		Fatalf("fn should success but failed with err: %v", err)
	}
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}