package main

import (
	"fmt"
	"time"

	"github.com/codyoss/retry"
)

func failingCodeGenerator() func() error {
	attempt := 0
	return func() error {
		attempt++
		return fmt.Errorf("I failed %d times 😢", attempt)
	}
}

func main() {
	// Create your own retry policy. It can be used across goroutines safely.
	b := &retry.ExponentialBackoff{
		Attempts:     5,
		InitialDelay: 0 * time.Millisecond,
	}
	failingCode := failingCodeGenerator()

	// final error will be returned if retries are exceeded
	err := retry.It(b, func() error {
		return failingCode()
	})
	fmt.Println(err)
	// Output: I failed 5 times 😢
}
