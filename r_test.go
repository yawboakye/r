package r

import (
	"errors"
	"testing"
	"time"

	"github.com/yawboakye/r/backoff"
)

func TestMaxRetries(t *testing.T) {
	fn := func(...interface{}) (interface{}, error) {
		return nil, errors.New("failed")
	}

	f := F{
		Fn:         fn,
		MaxRetries: 5,
		Backoff:    backoff.Exponential{time.Millisecond},
	}

	if _, err := f.Run(); err == nil {
		t.Fatal("expected non-nil error; got nil instead")
	}

	if f.MaxRetries != f.tries {
		t.Fatalf("expected maximum trials (%d); tried %d instead",
			f.MaxRetries, f.tries)
	}
}

func TestPassBeforeMaxRetries(t *testing.T) {
	f := F{
		Fn:         func(...interface{}) (interface{}, error) { return nil, nil },
		MaxRetries: 2,
		Backoff:    backoff.Exponential{time.Millisecond},
	}

	if _, err := f.Run(); err != nil {
		t.Fatalf("expected error to be nil; got=%v instead", err)
	}

	if f.tries == f.MaxRetries {
		t.Fatalf("expected fewer than max retries (%d); got=%d instead",
			f.MaxRetries, f.tries)
	}
}
