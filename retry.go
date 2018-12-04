// Package retry provides an simple api to retry function calls in a thread safe way.
package retry

import (
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
	Me = errors.New("Retry me")
)

// ExponentialBackoff holds the configuration of a backoff policy.
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
}

// It takes an ExponentialBackoff and a func that returns an err. If the function passed to this method returns an
// error it will be retired. It the number of attempts from the ExponentialBackoff is exceeded the final error
// the func returns will be returned.
//
// This function makes use of closures so any variables you would like to capture should be declared outside the
// invocation of this method.
func It(b *ExponentialBackoff, fn func() error) (err error) {
	b.mutex.Do(b.validate)

	delay := b.InitialDelay
	for i := 0; i < b.Attempts; i++ {
		if i != 0 {
			time.Sleep(delay)

			delay = time.Duration(float64(delay) * b.Factor)
			if delay > b.MaxDelay {
				delay = b.MaxDelay
			}
			if !b.skipJitter {
				delta := b.Jitter * float64(delay)
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

func (b *ExponentialBackoff) validate() {
	if b.Attempts < 1 {
		b.Attempts = 1
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
	rand.Seed(time.Now().Unix())
}
