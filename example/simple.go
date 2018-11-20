package main

import (
	"fmt"

	"github.com/codyoss/retry"
)

func squareIfInputIsThreeGenerator() func() int {
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
	squareIfInputIsThree := squareIfInputIsThreeGenerator()

	var result int
	retry.It(retry.DefaultExponentialBackoff, func() (err error) {
		result = squareIfInputIsThree()
		if result == 0 {
			return retry.Me
		}
		return
	})
	fmt.Println(result)
	//Output: 9
}
