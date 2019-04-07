package main

import (
	"context"
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
	backoff := &retry.Backoff{
		Attempts:     5,
		InitialDelay: 0 * time.Millisecond,
	}
	failingCode := failingCodeGenerator()

	// final error will be returned if retries are exceeded
	err := backoff.It(context.Background(), func(ctx context.Context) error {
		return failingCode()
	})
	fmt.Println(err)
	// Output: I failed 5 times 😢
}
