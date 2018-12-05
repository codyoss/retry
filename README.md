# retry

retry is a package that enables retrying code.

It has:

- a simple api(currently one public method)
- sane defaults
- validation logic to catch errors and optimize allocations
- a thread safe model
- a fun ergonomic syntax(at least I think soo :laughing:)

## Installation

```
go get github.com/codyoss/retry
```

## Examples

An example calling a function that does not return an error:

```go
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
```

An example calling a function that returns an error:

```go
var result int
retry.It(retry.DefaultExponentialBackoff, func() (err error) {
    result, err = squareOnThirdAttempt()
    return
})
fmt.Println(result)
// Output: 9
```

An example of what happens when retries are exceeded:

```go
// Create your own retry policy. It can be used across goroutines safely.
b := &retry.ExponentialBackoff{
    Attempts:     5,
    InitialDelay: 0 * time.Millisecond,
}
attempt := 0

// final error will be returned if retries are exceeded
err := retry.It(b, func() error {
    attempt++
    return fmt.Errorf("I failed %d times 😢", attempt)
})
fmt.Println(err)
// Output: I failed 5 times 😢
```

For full examples see [the examples folder](examples/)

## TODOs

- Add one more method that supports a passed in context
- take some feedback

## Disclaimer

Until this api hits v0.1.0 it might change a bit. After that I will try to keep things rather stable. After I have
gotten enough feedback I will v1.0.0 and make the same promise Go does. Cheers!