package r

import (
	"math/rand"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	i := NewInterval(time.Second, func() {})
	if cap(i.timer) != 1 {
		t.Fatalf("expected channel buffer size 1; got=%d instead", cap(i.timer))
	}

	if i.running {
		t.Fatalf("expected running to be false")
	}
}

func wait(secs int) {
	asDur := time.Duration(secs)
	time.Sleep(asDur * time.Second)
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

func TestStop(t *testing.T) {
	i := NewInterval(time.Second, func() {})
	i.Start()
	wait(1)
	i.Stop()

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
