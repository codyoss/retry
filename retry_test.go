package retry_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/codyoss/retry"
)

func TestIt(t *testing.T) {
	tests := []struct {
		name     string
		backoff  *retry.Backoff
		deadline time.Duration // 0 means no deadline
		setup    func() func(context.Context) (string, error)
		want     string
		wantErr  error
	}{
		{
			name:    "works",
			backoff: &retry.Backoff{Attempts: 3, Factor: 1.0, InitialDelay: 1 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1},
			setup: func() func(context.Context) (string, error) {
				attempt := 1
				return func(ctx context.Context) (string, error) {
					if attempt != 3 {
						attempt++
						return "", fmt.Errorf("Something is wrong: %d", attempt)
					}
					return "It Worked", nil
				}
			},
			want:    "It Worked",
			wantErr: nil,
		},
		{
			name:    "works, uses maxDelay",
			backoff: &retry.Backoff{Attempts: 3, Factor: 2.0, InitialDelay: 1 * time.Millisecond, MaxDelay: 2 * time.Millisecond, Jitter: .1},
			setup: func() func(context.Context) (string, error) {
				attempt := 1
				return func(ctx context.Context) (string, error) {
					if attempt != 3 {
						attempt++
						return "", fmt.Errorf("Something is wrong: %d", attempt)
					}
					return "It Worked", nil
				}
			},
			want:    "It Worked",
			wantErr: nil,
		},
		{
			name:    "not work, no retries",
			backoff: &retry.Backoff{Attempts: 1, Factor: 1.0, InitialDelay: 500 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1},
			setup: func() func(context.Context) (string, error) {
				return func(ctx context.Context) (string, error) {
					return "", retry.Me
				}
			},
			want:    "",
			wantErr: retry.Me,
		},
		{
			name:    "fast fail on context canceled",
			backoff: retry.ExponentialBackoff,
			setup: func() func(context.Context) (string, error) {
				return func(ctx context.Context) (string, error) {
					return "", context.Canceled
				}
			},
			want:    "",
			wantErr: context.Canceled,
		},
		{
			name:     "deadline exceeded during wait",
			backoff:  &retry.Backoff{Attempts: 5, InitialDelay: 50 * time.Millisecond, Factor: 1.0},
			deadline: 10 * time.Millisecond,
			setup: func() func(context.Context) (string, error) {
				return func(ctx context.Context) (string, error) {
					return "", retry.Me
				}
			},
			want:    "",
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.deadline > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tt.deadline)
				defer cancel()
			}

			fn := tt.setup()
			got, err := retry.It(ctx, tt.backoff, fn)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got err %v, want %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		backoff *retry.Backoff
		setup   func() func(context.Context) error
		wantErr error
	}{
		{
			name:    "works",
			backoff: &retry.Backoff{Attempts: 3, Factor: 1.0, InitialDelay: 1 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1},
			setup: func() func(context.Context) error {
				attempt := 1
				return func(ctx context.Context) error {
					if attempt != 3 {
						attempt++
						return fmt.Errorf("Something is wrong: %d", attempt)
					}
					return nil
				}
			},
			wantErr: nil,
		},
		{
			name:    "not work, no retries",
			backoff: &retry.Backoff{Attempts: 1, Factor: 1.0, InitialDelay: 500 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1},
			setup: func() func(context.Context) error {
				return func(ctx context.Context) error {
					return retry.Me
				}
			},
			wantErr: retry.Me,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := tt.setup()
			err := retry.Run(context.Background(), tt.backoff, fn)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got err %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestBackoffRun(t *testing.T) {
	tests := []struct {
		name    string
		backoff *retry.Backoff
		setup   func() func(context.Context) error
		wantErr error
	}{
		{
			name:    "works",
			backoff: &retry.Backoff{Attempts: 3, Factor: 1.0, InitialDelay: 1 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1},
			setup: func() func(context.Context) error {
				attempt := 1
				return func(ctx context.Context) error {
					if attempt != 3 {
						attempt++
						return fmt.Errorf("Something is wrong: %d", attempt)
					}
					return nil
				}
			},
			wantErr: nil,
		},
		{
			name:    "not work, no retries",
			backoff: &retry.Backoff{Attempts: 1, Factor: 1.0, InitialDelay: 500 * time.Millisecond, MaxDelay: retry.Forever, Jitter: .1},
			setup: func() func(context.Context) error {
				return func(ctx context.Context) error {
					return retry.Me
				}
			},
			wantErr: retry.Me,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := tt.setup()
			err := tt.backoff.Run(context.Background(), fn)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("got err %v, want %v", err, tt.wantErr)
			}
		})
	}
}
