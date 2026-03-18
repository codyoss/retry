# retry

retry is a package that enables retrying code.

[![GoDoc](https://godoc.org/github.com/codyoss/retry?status.svg)](https://godoc.org/github.com/codyoss/retry)
[![Build Status](https://github.com/codyoss/retry/actions/workflows/ci.yml/badge.svg)](https://github.com/codyoss/retry/actions/workflows/ci.yml)
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
result, err := retry.It(context.Background(), retry.ExponentialBackoff, func(ctx context.Context) (int, error) {
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
//Output: 9
```

// Removed the backoff method example since it was removed from the API to support generics.

An example calling a function that returns an error:

```go
result, err := retry.It(context.Background(), retry.ExponentialBackoff, func(ctx context.Context) (int, error) {
    return squareOnThirdAttempt()
})
if err != nil {
    // TODO: handle error
}
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
err := backoff.Run(context.Background(), func(ctx context.Context) error {
    attempt++
    return fmt.Errorf("I failed %d times \U0001f622", attempt)
})
if err != nil {
    // TODO: handle error
}
// Output: I failed 5 times 😢
```

retry also supports calling code that is context aware!

```go
noopFn := func(ctx context.Context) {}

// can you use the context to set a max amount time to retry for
ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

// the context passed in is forwarded to the function provided
err := retry.Run(ctx, retry.ConstantDelay, func(ctx context.Context) error {
    noopFn(ctx)
    return retry.Me
})

fmt.Printf("%v\n", err)
// Output: context deadline exceeded
```

For full examples with more documentation see [the examples folder](examples/)
