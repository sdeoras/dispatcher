package dispatcher

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestDispatcher(t *testing.T) {
	counter := new(int32)

	d := New(5) // a 5 goroutine dispatcher

	for i := 0; i < 10; i++ {
		f := func() {
			time.Sleep(time.Second)
			atomic.AddInt32(counter, 1)
		}
		d.Do(f)
	}

	time.Sleep(time.Millisecond * 20) // sleep a bit just so that functions get scheduled

	for d.IsRunning() {
		time.Sleep(time.Second)
	}

	if *counter != 10 {
		t.Fatal("not all functions executed. expected 10, got:", *counter)
	}
}
