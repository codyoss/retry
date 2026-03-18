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

	result, err := retry.It(context.Background(), retry.ExponentialBackoff, func(ctx context.Context) (int, error) {
		// Put code you would like to retry here. If you return an error and have not exceeded retries the code in
		// in this block will be executed again based on the backoff policy provided.
		res := squareOnThirdAttempt()
		if res == 0 {
			return 0, retry.Me
		}
		return res, nil
	})
	if err != nil {
		// TODO: handle error
	}
	fmt.Println(result)
	// Output: 9
}
