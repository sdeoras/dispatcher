package dispatcher

import (
	"sync"
	"sync/atomic"
	"time"
)

// Dispatcher defines an interface for a function dispatcher
type Dispatcher interface {
	// Do causes function to be scheduled for execution.
	// Execution trigger is implementation specific.
	Do(f func())
	// IsRunning provides activity status of the dispatcher
	IsRunning() bool
}

// dispatcher implements Dispatcher interface such that is allows user
// to limit the number of active goroutines running at a time.
type dispatcher struct {
	funcStack     *queue
	maxDispatched int32
	active        *int32
	poke          chan struct{}
	mu            sync.Mutex
}

// New provides a new instance of dispatcher
func New(numConcurrent int32) *dispatcher {
	d := new(dispatcher)
	d.funcStack = new(queue)
	d.maxDispatched = numConcurrent
	d.active = new(int32)
	d.poke = make(chan struct{})
	d.start() // starts a daemon that will schedule pending funcs
	return d
}

// queue is a list of func
type queue []func()

func (s *queue) len() int {
	return len(*s)
}

// push pushes new entry at the end of the list
func (s *queue) push(f func()) {
	*s = append(*s, f)
}

// pop pulls from the front of the queue
func (s *queue) pop() func() {
	if len(*s) > 0 {
		f := (*s)[0]
		*s = (*s)[1:]
		return f
	}

	return nil
}

func (d *dispatcher) IsRunning() bool {
	return *(d.active) > 0
}

func (d *dispatcher) pending() int {
	return d.funcStack.len()
}

func (d *dispatcher) Do(f func()) {
	// lock
	d.mu.Lock()
	defer d.mu.Unlock()

	// push into queue
	d.funcStack.push(f)

	d.dispatch()
}

// dispatch is an internal function
func (d *dispatcher) dispatch() {
	for *(d.active) < d.maxDispatched {
		f := d.funcStack.pop()
		if f == nil {
			break
		}

		// increment the active counter
		atomic.AddInt32(d.active, 1)

		go func(active *int32, poke chan struct{}) {
			f()
			atomic.AddInt32(active, -1)
			d.poke <- struct{}{}
		}(d.active, d.poke)
	}
}

// start is an internal function that monitors the active functions and dispatches new from pending list
func (d *dispatcher) start() {
	go func() {
		// run infinite loop waiting every second
		for {
			d.mu.Lock()
			d.dispatch()
			d.mu.Unlock()
			select {
			case <-d.poke:
			case <-time.After(time.Second):
			}
		}
	}()
}
