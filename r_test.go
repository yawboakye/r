package r

import (
	"errors"
	"testing"
	"time"

	"github.com/yawboakye/r/backoff"
)

func TestMaxRetries(t *testing.T) {
	fn := func(...interface{}) error {
		return errors.New("failed")
	}

	f := F{
		Fn:         fn,
		MaxRetries: 5,
		Backoff:    backoff.Exponential{time.Millisecond},
	}

	f.Run()

	if f.Err() == nil {
		t.Fatal("expected non-nil error; got nil instead")
	}

	if f.MaxRetries != f.tries {
		t.Fatalf("expected maximum trials (%d); tried %d instead",
			f.MaxRetries, f.tries)
	}
}

func TestPassBeforeMaxRetries(t *testing.T) {
	f := F{
		Fn:         func(...interface{}) error { return nil },
		MaxRetries: 2,
		Backoff:    backoff.Exponential{time.Millisecond},
	}

	f.Run()
	err := f.Err()
	if err != nil {
		t.Fatalf("expected error to be nil; got=%q instead", err)
	}

	if f.tries == f.MaxRetries {
		t.Fatalf("expected fewer than max retries (%d); got=%d instead",
			f.MaxRetries, f.tries)
	}
}
