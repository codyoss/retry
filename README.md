# retry

retry is a package that enables retrying code.

It has:

- a simple api(currently one public method)
- sane defaults
- validation logic to catch errors and optimize allocations
- a thread safe model
- a fun syntax(at least I think soo :laughing:)

## Examples

An example calling a function that does not return an error:

```go
package main

import (
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
    retry.It(retry.DefaultExponentialBackoff, func() (err error) {
        result = squareOnThirdAttempt()
        if result == 0 {
            return retry.Me
        }
        return
    })
    fmt.Println(result)
    //Output: 9
}
```

An example calling a function that returns an error

```go
package main

import (
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
    retry.It(retry.DefaultExponentialBackoff, func() (err error) {
        result, err = squareOnThirdAttempt()
        return
    })
    fmt.Println(result)
    // Output: 9
}
```

## TODOs

- Add one more method that supports a passed in context
- take some feedback

## Disclaimer

Until this api hits v0.1.0 it might change a bit. After that I will try to keep things rather stable. After I have
gotten enough feedback I will v1.0.0 and make the same promise Go does. Cheers!