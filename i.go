package r

import (
	"errors"
	"time"
)

// BEWARE: This is an overengineered solution to
// a tiny problem. If you want to run a function
// at any interval (I suppose using the sleep pattern),
// this is all you need:
//
// for {
//   go func() {
//	   statements...
//   }
//   time.Sleep(duration)
// }

// Trying to stop an interval that hasn't already
// been running returns this error.
var stopErr = errors.New("i: interval not running")

// I is a manifest for a function/method that should
// be run at the given interval. The next call is
// schedule right at the beginning of the current call
// so that the pulse isn't affected by how long the
// function takes to run.
type I struct {
	fn      func()
	intv    time.Duration
	timer   chan struct{}
	running bool
}

// NewInterval returns an initialized interval ready
// to be started. In most cases, if not all, this is
// how you'd want to create a new interval.
func NewInterval(intv time.Duration, f func()) *I {
	return &I{
		fn:      f,
		intv:    intv,
		timer:   make(chan struct{}, 1),
		running: false,
	}
}

// Starts the intervaled execution of the function.
// Returns true if the interval was started, false otherwise.
func (i *I) Start() bool {
	if i.running {
		return false
	}

	i.running = true
	// Start a goroutine to control pulse of execution.
	go func() {
		for i.running {
			i.timer <- struct{}{}
			time.Sleep(i.intv)
		}
		close(i.timer)
	}()

	go func() {
		for {
			// We continue to start new goroutines to run the function
			// until the interval is stopped (aka timer channel is closed)
			// in which case we'd receive a non-OK default value.
			_, ok := <-i.timer
			if !ok {
				break
			}

			go i.fn()
		}
	}()

	return true
}

// Stop stops a running interval.
// If the interval was already stopped or not running,
// an error is returned.
func (i *I) Stop() error {
	if !i.running {
		return stopErr
	}

	i.running = false
	return nil
}
