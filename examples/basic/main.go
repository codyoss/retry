package main

import (
	"context"
	"fmt"

	"github.com/codyoss/retry"
)

func squareOnThirdAttemptGenerator() func() int {
	attempt := 1
	return func() int {
		if attempt != 3 {
			attempt++
			return 0
		}
		return attempt * attempt
	}
}

func main() {
	squareOnThirdAttempt := squareOnThirdAttemptGenerator()

	var result int
	retry.It(context.Background(), retry.ExponentialBackoff, func(ctx context.Context) (err error) {
		// Put code you would like to retry here. If you return an error and have not exceeded retries the code in
		// in this block will be executed again based on the backoff policy provided.
		result = squareOnThirdAttempt()
		if result == 0 {
			return retry.Me
		}
		return
	})
	fmt.Println(result)
	// Output: 9
}
