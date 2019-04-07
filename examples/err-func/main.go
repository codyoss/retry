package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/codyoss/retry"
)

func squareOnThirdAttemptGenerator() func() (int, error) {
	attempt := 1
	return func() (int, error) {
		if attempt != 3 {
			attempt++
			return 0, errors.New("uh oh")
		}
		return attempt * attempt, nil
	}
}

func main() {
	squareOnThirdAttempt := squareOnThirdAttemptGenerator()

	var result int
	retry.It(context.Background(), retry.ConstantDelay, func(ctx context.Context) (err error) {
		result, err = squareOnThirdAttempt()
		return
	})
	fmt.Println(result)
	// Output: 9
}
