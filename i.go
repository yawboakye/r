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
var errStop = errors.New("i: interval not running")

// Trying to abort an already aborted interval returns
// this error.
var errAborted = errors.New("i: interval already aborted")

// I is a manifest for a function/method that should
// be run at the given interval. The next call is
// schedule right at the beginning of the current call
// so that the pulse isn't affected by how long the
// function takes to run.
type I struct {
	fn      func()
	intv    time.Duration
	timer   chan struct{}
	abort   chan struct{}
	done    chan struct{}
	running bool
	aborted bool
	started int
}

// NewInterval returns an initialized interval ready
// to be started. In most cases, if not all, this is
// how you'd want to create a new interval.
func NewInterval(intv time.Duration, f func()) *I {
	return &I{
		fn:      f,
		intv:    intv,
		timer:   make(chan struct{}),
		abort:   make(chan struct{}),
		done:    make(chan struct{}),
		running: false,
	}
}

// Start starts the intervaled execution of the function.
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

			go func() {
				i.fn()
				i.done <- struct{}{}
			}()
			i.started++

			select {
			case <-i.abort:
				// Exit the parent goroutine if we receive
				// on the `abort` channel. This should close
				// all running tasks.
				return

			case <-i.done:
				continue
			}
		}
	}()

	return true
}

// Started returns a count of calls that have been started
// so far. Even after calling Stop or Abort, this value is
// still available.
func (i *I) Started() int { return i.started }

// Stop stops a running interval.
// If the interval was already stopped or not running,
// an error is returned. It doesn't stop calls that
// have already been started. It just stops the interval
// from starting new ones. If you want to abort running
// functions then Abort instead.
func (i *I) Stop() error {
	if !i.running {
		return errStop
	}

	close(i.timer)
	i.running = false
	return nil
}

// Abort aborts all running functions. After a call
// to Stop, functions that had already been started
// are not terminated. If you want to stop the
// interval and abort all functions that had been
// started but not completed, abort instead.
func (i *I) Abort() error {
	if i.aborted {
		return errAborted
	}

	close(i.timer)
	close(i.abort)
	i.running = false
	i.aborted = true
	return nil
}
