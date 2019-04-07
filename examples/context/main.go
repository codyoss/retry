package main

import (
	"context"
	"fmt"
	"time"

	"github.com/codyoss/retry"
)

func main() {
	noopFn := func(ctx context.Context) {}

	// can you use the context to set a max amount time to retry for
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := retry.It(ctx, retry.ConstantDelay, func(ctx context.Context) (err error) {
		// The context that is provided to `retry.ItContext` is passed into this function so you may call your context
		// aware code and have it tied to your parent context
		noopFn(ctx)
		return retry.Me
	})

	// In this case the context causes the retrying to timeout, so that is the error that is reported back to the
	// calling function.
	fmt.Printf("%v\n", err)
	// Output: context deadline exceeded
}
