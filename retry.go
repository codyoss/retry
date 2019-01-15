// Package retry provides an simple api to retry function calls in a thread safe way.
package retry

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"
)

const (
	// Forever is a shorthand for a time.Duration of one year. This value can be used as an input to ExponentialBackoff
	// for the MaxDelay field should the user not want to ever hit that limit.
	Forever = 365 * 24 * 60 * time.Minute
)

var (
	// DefaultConstantDelay is a backoff policy that will attempt a call 5 times with 500 milliseconds between
	// subsequent calls.
	DefaultConstantDelay = &ExponentialBackoff{
		Attempts:     5,
		InitialDelay: 500 * time.Millisecond,
		Factor:       1,
	}
	// DefaultExponentialBackoff provides a sane exponential backoff policy.
	DefaultExponentialBackoff = &ExponentialBackoff{
		Attempts:     5,
		InitialDelay: 500 * time.Millisecond,
		MaxDelay:     8 * time.Second,
		Factor:       2.0,
		Jitter:       0.1,
	}

	// Me is an error that can be returned in from a function. To be used in the function passed to It if no other
	// error makes sense or if you don't care to return that actual error. This is error variable is simply sugar.
	Me = errors.New("retry me")
)

// ExponentialBackoff holds the configuration of a backoff policy. Once values are set and this backoff is used any
// modification to this struct will not affect behavior. Fields get frozen to un-exported variables upon first use.
type ExponentialBackoff struct {
	// Attempts it the max number of times a function will be retried. Will always be treaded as a value >= 1.
	Attempts int
	// InitialDelay is the starting delay should the first attempt fail.
	InitialDelay time.Duration
	// MaxDelay is the max amount of time to try in between retry attempts.
	MaxDelay time.Duration
	// Factor is the amount that will be multiplied to the previous delay to calculate the next delay. This value
	// will always be treaded as a value >= 1.0. A value of 2.0 would be a standard exponential backoff.
	Factor float64
	// Jitter is a way to add a bit of randomness into your delay. Setting this value helps avoid what is known as the
	// thundering herd problem. For example if a value of .1 is set and your delay is 500 milliseconds the Jitter would
	// tranform that value into a number between 490 and 510 milliseconds.
	Jitter float64

	mutex      sync.Once
	skipJitter bool

	// Frozen fields
	attempts     int
	initialDelay time.Duration
	maxDelay     time.Duration
	factor       float64
	jitter       float64
}

// It takes an ExponentialBackoff and a func that returns an err. If the function passed to this method returns an
// error it will be retired. It the number of attempts from the ExponentialBackoff is exceeded the final error
// the func returns will be returned.
//
// This function makes use of closures so any variables you would like to capture should be declared outside the
// invocation of this method.
func It(b *ExponentialBackoff, fn func() error) (err error) {
	b.mutex.Do(b.validateAndFreeze)

	delay := b.initialDelay
	for i := 0; i < b.attempts; i++ {
		if i != 0 {
			time.Sleep(delay)

			delay = time.Duration(float64(delay) * b.factor)
			if delay > b.maxDelay {
				delay = b.maxDelay
			}
			if !b.skipJitter {
				delta := b.jitter * float64(delay)
				minDelay := float64(delay) - delta
				maxDelay := float64(delay) + delta
				delay = time.Duration(minDelay + (rand.Float64() * (maxDelay - minDelay + 1)))
			}
		}
		err = fn()
		if err == nil {
			return
		}
	}
	return
}

// ItContext is a the same as `It` but context aware. This methods can be used to set an overall timeout. It will also
// pass the provided context to the function provided. Thus, any code you call within the retry block can share the
// same parent context.
func ItContext(ctx context.Context, b *ExponentialBackoff, fn func(context.Context) error) (err error) {
	b.mutex.Do(b.validateAndFreeze)

	delay := b.initialDelay
	for i := 0; i < b.attempts; i++ {
		if i != 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}

			delay = time.Duration(float64(delay) * b.factor)
			if delay > b.maxDelay {
				delay = b.maxDelay
			}
			if !b.skipJitter {
				delta := b.jitter * float64(delay)
				minDelay := float64(delay) - delta
				maxDelay := float64(delay) + delta
				delay = time.Duration(minDelay + (rand.Float64() * (maxDelay - minDelay + 1)))
			}
		}
		err = fn(ctx)
		if err == nil {
			return
		}
	}
	return
}

func (b *ExponentialBackoff) validateAndFreeze() {
	// validation
	if b.Attempts < 1 {
		b.Attempts = 1
	}
	if b.InitialDelay < 0 {
		b.InitialDelay = 0
	}
	if b.MaxDelay == 0 {
		b.MaxDelay = Forever
	}
	if b.Factor < 1 {
		b.Factor = 1
	}
	if b.Jitter < 0 {
		b.Jitter = 0
	}
	if b.Jitter == 0 {
		b.skipJitter = true
	}

	// freeze
	rand.Seed(time.Now().Unix())
	b.attempts = b.Attempts
	b.initialDelay = b.InitialDelay
	b.maxDelay = b.MaxDelay
	b.factor = b.Factor
	b.jitter = b.Jitter
}
