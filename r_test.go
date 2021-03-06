package r

import (
	"errors"
	"testing"
	"time"

	"github.com/yawboakye/r/backoff"
)

var (
	err    = errors.New("error")
	failFn = func(...interface{}) (interface{}, error) { return nil, err }
	passFn = func(...interface{}) (interface{}, error) { return nil, nil }
)

func TestMaxRetries(t *testing.T) {
	f := F{
		Fn:         failFn,
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
		Fn:         passFn,
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

func TestSingleUse(t *testing.T) {
	f := F{
		Fn:         passFn,
		MaxRetries: 2,
		Backoff:    backoff.Linear{time.Millisecond},
	}

	f.Run() // Run, without caring about the returned value
	if !f.used {
		t.Fatal("expected f.used to be true; got false instead")
	}

	tries := f.tries
	_, err := f.Run()
	if f.tries != tries || err == nil {
		t.Fatal("expected no trials. but a trial happened")
	}
}

func TestPassArgs(t *testing.T) {
	err := errors.New("div: division by zero")
	div := func(a, b int) (int, error) {
		if b == 0 {
			return 0, err
		}
		return a / b, nil
	}

	f := F{
		MaxRetries: 2,
		Backoff:    backoff.Linear{time.Millisecond},
		Fn: func(a ...interface{}) (interface{}, error) {
			return div(a[0].(int), a[1].(int))
		},
	}

	_, e := f.Run(0, 0)
	if e != err {
		t.Fatalf("expected %v; got=%v instead", err, e)
	}
}
