# retry

retry is a a package that enables retrying code. 

It has:
- a simple api(currently one public method) 
- sane defaults
- validation logic to catch errors and optimize allocations
- a thread safe model
- a fun syntax(at least I think soo :laughing:)

## Examples

```go
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
```

## TODOs

- Add one more method that supports a passed in context
- take some feedback

## Disclaimer

Until this api hits v0.1.0 it might change a bit. After that I will try to keep things rather stable. After I have
gotten enough feedback I will v1.0.0 and make the same promise Go does. Cheers!