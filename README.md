# retry

retry is a package that enables retrying code.

It has:

- a simple api, only two public functions
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

### retry.ItContext

retry also supports calling code that is context aware via `retry.ItContext`

```go
noopFn := func(ctx context.Context) {}

// can you use the context to set a max amount time to retry for
ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
defer cancel()

err := retry.ItContext(ctx, retry.DefaultConstantDelay, func(ctx context.Context) (err error) {
    noopFn(ctx)
    return retry.Me
})

fmt.Printf("%v\n", err)
// Output: context deadline exceeded
```

For full examples with more documentation see [the examples folder](examples/)