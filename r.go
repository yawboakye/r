package r

import (
	"time"

	"github.com/yawboakye/r/backoff"
)

// An F is a manifest for a retry-able function/method.
// It defines what function to retry, maximum retries
// before returning the ensuing error, and the backoff
// strategy to use between retries.
type F struct {
	Fn         func(...interface{}) (interface{}, error)
	MaxRetries int
	Backoff    backoff.Strategy
	tries      int
}

// Returns the number of times the function was tried.
// This is only informational if the function succeeded
// after a number of calls. In that case it will be
// different and lower than MaxRetries.
func (f *F) Tried() int { return f.tries }

func (f *F) exhausted() bool { return f.tries == f.MaxRetries }

// Run runs the function, retrying on failure until the
// maximum number of retries is exceeded.
func (f *F) Run(args ...interface{}) (interface{}, error) {
	for {
		f.tries++

		res, err := f.Fn(args)
		if err == nil || f.exhausted() {
			return res, err
		}

		// The bad thing happened.
		// Wait for the duration decided by the backoff
		// strategy, and then try again.
		f.wait()
		continue
	}
}

// wait waits for a period between two retries
// of a function. How long it waits for depends
// on the backoff strategy.
func (f *F) wait() {
	time.Sleep(f.Backoff.WaitDur(f.tries))
}
