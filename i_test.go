package r

import (
	"math/rand"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	i := NewInterval(time.Second, func() {})
	if cap(i.timer) != 0 {
		t.Fatalf("expected channel buffer size 1; got=%d instead", cap(i.timer))
	}

	if i.running {
		t.Fatalf("expected running to be false")
	}
}

// TODO(yawboakye): Revisit this test case. I'm not
// confident that I've tested it the best way possible.
func TestStart(t *testing.T) {
	var calls int
	var expected = 1 + rand.Intn(5)

	i := NewInterval(time.Second, func() { calls++ })

	// Start the interval, wait for the expected number of
	// calls and then stop it.
	i.Start()
	wait(expected)
	i.running = false // instead of `i.Stop()`

	// expect the function to be have been called
	// the expected number of times.
	if calls != expected {
		t.Fatalf("expected=%d; got=%d instead", expected, calls)
	}
}

func TestStartStart(t *testing.T) {
	i := NewInterval(time.Second, func() {})
	i.Start()
	if ok := i.Start(); ok {
		t.Fatalf("started an already running interval")
	}
	i.Stop()
}

func TestStarted(t *testing.T) {
	i := NewInterval(time.Second, func() {})
	i.Start()
	wait(1)

	if i.Started() != i.started {
		t.Fatalf("expected=%d; got=%d instead", i.started, i.Started())
	}
}

func TestStop(t *testing.T) {
	var j int
	i := NewInterval(2*time.Second, func() { j++ })
	i.Start()
	wait(3)
	i.Stop()

	if j != i.started {
		t.Fatalf("expected completions=%d; got=%d instead", i.started, j)
	}

	// expect `running` to be false
	if i.running {
		t.Fatalf("interval still running. not stopped")
	}

	// expect the timer channel to be closed
	if _, ok := <-i.timer; ok {
		t.Fatalf("expected timer channel to be closed")
	}

	// stopping an already stopped interval returns error
	if err := i.Stop(); err == nil {
		t.Fatalf("expected=%v; got=%v instead", errStop, nil)
	}
}

func TestAbort(t *testing.T) {
	var j int

	// slowFn waits for 2 seconds and then increments j.
	slowFn := func() {
		time.Sleep(2 * time.Second)
		j++
	}

	i := NewInterval(time.Second, slowFn)
	i.Start()
	wait(3)
	i.Abort()

	if j >= i.started {
		t.Fatalf("expected completions=%d; got=%d instead", j, i.started)
	}

	// expect `aborted` to be true
	if !i.aborted {
		t.Fatal("interval not aborted")
	}

	// expect the abort channel to be closed
	if _, ok := <-i.abort; ok {
		t.Fatal("expected abort channel to be closed")
	}

	// expect `running` to be false
	if i.running {
		t.Fatal("interval still running. not stopped")
	}

	// expect the timer channel to be closed
	if _, ok := <-i.timer; ok {
		t.Fatal("expected timer channel to be closed")
	}

	// aborting an already aborted interval returns error
	if err := i.Abort(); err == nil {
		t.Fatalf("expected=%v; got=%v instead", errAborted, nil)
	}
}

func wait(t int) { time.Sleep(time.Duration(t) * time.Second) }
