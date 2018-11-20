package retry_test

import (
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
		backoff *retry.ExponentialBackoff
		want    string
		wantErr error
	}{
		{"works", &retry.ExponentialBackoff{Attempts: 3, Factor: 1.0, InitialDelay: 1 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1}, "It Worked", nil},
		{"not work, no retries", &retry.ExponentialBackoff{Attempts: 1, Factor: 1.0, InitialDelay: 500 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1}, "", retry.Me},
		{"not work, one retry", &retry.ExponentialBackoff{Attempts: 2, Factor: 1.0, InitialDelay: 1 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1}, "", retry.Me},
	}

	for _, tt := range tests {
		fn := aFunctionThatWorksAfterCalledThreeTimes()
		t.Run(tt.name, func(t *testing.T) {
			var got string
			err := retry.It(tt.backoff, func() (err error) {
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
