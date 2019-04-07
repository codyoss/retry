# retry

retry is a package that enables retrying code.

[![GoDoc](https://godoc.org/github.com/codyoss/retry?status.svg)](https://godoc.org/github.com/codyoss/retry)
[![Build Status](https://cloud.drone.io/api/badges/codyoss/retry/status.svg)](https://cloud.drone.io/codyoss/retry)
[![codecov](https://codecov.io/gh/codyoss/retry/branch/master/graph/badge.svg)](https://codecov.io/gh/codyoss/retry)
[![Go Report Card](https://goreportcard.com/badge/github.com/codyoss/retry)](https://goreportcard.com/report/github.com/codyoss/retry)

It has:

- a simple api, only one public function!
- sane defaults
- validation logic to catch errors and optimize allocations
- a thread safe model
- support for `context.Context`
- a fun ergonomic syntax(at least I think soo :laughing:)

## Installation

```bash
go get github.com/codyoss/retry
```

## Examples

An example calling a function that does not return an error:

```go
var result int
retry.It(context.Background(), retry.ExponentialBackoff, func(ctx context.Context) (err error) {
    result = squareOnThirdAttempt()
    if result == 0 {
        return retry.Me
    }
    return
})
fmt.Println(result)
//Output: 9
```

An alternate syntax that does the same thing as above:

```go
var result int
// You could also create your own backoff policy
backoff := retry.ExponentialBackoff
// This just calls the package level It function under the hood.
backoff.It(context.Background(), func(ctx context.Context) (err error) {
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
retry.It(context.Background(), retry.ExponentialBackoff, func(ctx context.Context) (err error) {
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
err := retry.It(context.Background(), b, func(ctx context.Context) error {
    attempt++
    return fmt.Errorf("I failed %d times 😢", attempt)
})
fmt.Println(err)
// Output: I failed 5 times 😢
```

retry also supports calling code that is context aware!

```go
noopFn := func(ctx context.Context) {}

// can you use the context to set a max amount time to retry for
ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

// the context passed in is forwarded to the function provided
err := retry.It(ctx, retry.ConstantDelay, func(ctx context.Context) (err error) {
    noopFn(ctx)
    return retry.Me
})

fmt.Printf("%v\n", err)
// Output: context deadline exceeded
```

For full examples with more documentation see [the examples folder](examples/)

## Disclaimer

Until this api hits v0.1.0 it might change a bit. After that I will try to keep things rather stable. After I have gotten enough feedback I will v1.0.0 and make the same promise Go does. Cheers!