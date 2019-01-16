package retry

import (
	"testing"
	"time"
)

func TestFreezeBackoffAfterFirstUse(t *testing.T) {
	attempts := 1
	initialDelay := 1 * time.Millisecond
	maxDelay := 1 * time.Millisecond
	factor := 1.0
	jitter := 1.0

	b := &Backoff{
		Attempts:     attempts,
		InitialDelay: initialDelay,
		MaxDelay:     maxDelay,
		Factor:       factor,
		Jitter:       jitter,
	}

	It(b, func() error {
		return nil
	})

	b.Attempts = 7
	b.InitialDelay = 7 * time.Millisecond
	b.MaxDelay = 7 * time.Millisecond
	b.Factor = 7.0
	b.Jitter = 7.0

	It(b, func() error {
		return nil
	})

	if b.Attempts != 7 || b.attempts != attempts ||
		b.InitialDelay != 7*time.Millisecond || b.initialDelay != initialDelay ||
		b.MaxDelay != 7*time.Millisecond || b.maxDelay != maxDelay ||
		b.Factor != 7.0 || b.factor != factor ||
		b.Jitter != 7.0 || b.jitter != jitter {
		t.Error("public fields should change, private should have been frozen")
	}
}
