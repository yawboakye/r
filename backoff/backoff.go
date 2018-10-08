package backoff

import "time"

// Strategy is a backoff strategy.
type Strategy interface {
	// WaitDur determines how long to wait between
	// retries, based on the number of tries made so far.
	WaitDur(tries int) time.Duration
}

// Exponential is a simplified and predictable
// exponential backoff mechanism. See
// https://en.wikipedia.org/wiki/Exponential_backoff
// for more information.
type Exponential struct{ Dur time.Duration }

// A Linear backoff waits a constant duration between
// retries. Unlike an exponential backoff, feedback (i.e.
// number of tries so far) isn't used as information to
// determine how long to wait before the next try.
type Linear struct{ Dur time.Duration }

func (e Exponential) WaitDur(t int) time.Duration { return e.Dur << uint(t) }
func (l Linear) WaitDur(t int) time.Duration      { return l.Dur }
