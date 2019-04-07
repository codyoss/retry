package retry_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/codyoss/retry"
)

func aFunctionThatWorksAfterCalledThreeTimes() func() (string, error) {
	attempt := 1
	return func() (string, error) {
		if attempt != 3 {
			attempt++
			return "", fmt.Errorf("Something is wrong: %d", attempt)
		}
		res := "It Worked"
		return res, nil
	}
}

func TestIt(t *testing.T) {
	tests := []struct {
		name    string
		backoff *retry.Backoff
		want    string
		wantErr error
	}{
		{"works", &retry.Backoff{Attempts: 3, Factor: 1.0, InitialDelay: 1 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1}, "It Worked", nil},
		{"works, uses maxDelay", &retry.Backoff{Attempts: 3, Factor: 2.0, InitialDelay: 1 * time.Millisecond, MaxDelay: 2 * time.Millisecond, Jitter: .1}, "It Worked", nil},
		{"not work, no retries", &retry.Backoff{Attempts: 1, Factor: 1.0, InitialDelay: 500 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1}, "", retry.Me},
		{"not work, one retry", &retry.Backoff{Attempts: 2, Factor: 1.0, InitialDelay: 1 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1}, "", retry.Me},
	}

	for _, tt := range tests {
		fn := aFunctionThatWorksAfterCalledThreeTimes()
		t.Run(tt.name, func(t *testing.T) {
			var got string
			err := retry.It(context.Background(), tt.backoff, func(ctx context.Context) (err error) {
				got, err = fn()
				if err != nil {
					return retry.Me
				}
				return
			})

			if err != tt.wantErr {
				t.Error("Expected a nil error")
			}

			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestItWithContext(t *testing.T) {
	tests := []struct {
		name     string
		backoff  *retry.Backoff
		deadline time.Duration
		want     string
		wantErr  error
	}{
		{"works", &retry.Backoff{Attempts: 3, Factor: 1.0, InitialDelay: 1 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1}, 10 * time.Millisecond, "It Worked", nil},
		{"works, uses maxDelay", &retry.Backoff{Attempts: 3, Factor: 2.0, InitialDelay: 1 * time.Millisecond, MaxDelay: 2 * time.Millisecond, Jitter: .1}, 10 * time.Millisecond, "It Worked", nil},
		{"not work, no retries", &retry.Backoff{Attempts: 1, Factor: 1.0, InitialDelay: 500 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1}, 10 * time.Millisecond, "", retry.Me},
		{"not work, one retry", &retry.Backoff{Attempts: 2, Factor: 1.0, InitialDelay: 1 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1}, 10 * time.Millisecond, "", retry.Me},
		{"not work, deadline exceeded", retry.ExponentialBackoff, 10 * time.Millisecond, "", context.DeadlineExceeded},
	}

	for _, tt := range tests {
		fn := aFunctionThatWorksAfterCalledThreeTimes()
		t.Run(tt.name, func(t *testing.T) {
			var got string
			ctx, cancel := context.WithTimeout(context.Background(), tt.deadline)
			defer cancel()
			err := retry.It(ctx, tt.backoff, func(ctx context.Context) (err error) {
				got, err = fn()
				if err != nil {
					return retry.Me
				}
				return
			})

			if err != tt.wantErr {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}

			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}

func TestItWhenContextCanceled(t *testing.T) {
	fn := func(ctx context.Context) { return }

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	err := retry.It(ctx, retry.ExponentialBackoff, func(ctx context.Context) (err error) {
		fn(ctx)
		cancel()
		return retry.Me
	})

	if err != context.Canceled {
		t.Errorf("got %v, wanted a context.Canceled", err)
	}
}

func TestBackoffIt(t *testing.T) {
	tests := []struct {
		name    string
		backoff *retry.Backoff
		want    string
		wantErr error
	}{
		{"works", &retry.Backoff{Attempts: 3, Factor: 1.0, InitialDelay: 1 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1}, "It Worked", nil},
		{"works, uses maxDelay", &retry.Backoff{Attempts: 3, Factor: 2.0, InitialDelay: 1 * time.Millisecond, MaxDelay: 2 * time.Millisecond, Jitter: .1}, "It Worked", nil},
		{"not work, no retries", &retry.Backoff{Attempts: 1, Factor: 1.0, InitialDelay: 500 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1}, "", retry.Me},
		{"not work, one retry", &retry.Backoff{Attempts: 2, Factor: 1.0, InitialDelay: 1 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1}, "", retry.Me},
	}

	for _, tt := range tests {
		fn := aFunctionThatWorksAfterCalledThreeTimes()
		t.Run(tt.name, func(t *testing.T) {
			var got string
			err := tt.backoff.It(context.Background(), func(ctx context.Context) (err error) {
				got, err = fn()
				if err != nil {
					return retry.Me
				}
				return
			})

			if err != tt.wantErr {
				t.Error("Expected a nil error")
			}

			if got != tt.want {
				t.Errorf("got %s, want %s", got, tt.want)
			}
		})
	}
}
