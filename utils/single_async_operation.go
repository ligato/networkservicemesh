package utils

import (
	"sync/atomic"
)

const (
	notScheduled        = int32(0)
	running             = int32(1)
	scheduledAndRunning = int32(3)
)

//Operation means an object which can be executed
type Operation interface {
	//Run executes operation
	Run()
}

type singleAsyncOperation struct {
	body  func()
	state int32
}

//NewSingleAsyncOperation creates an operation which should be invoked once by run period. Can be used in cases where required the last run.
func NewSingleAsyncOperation(body func()) Operation {
	if body == nil {
		panic("body can not be nil")
	}
	return &singleAsyncOperation{body: body, state: notScheduled}
}

func (o *singleAsyncOperation) Run() {
	if !atomic.CompareAndSwapInt32(&o.state, notScheduled, running) {
		if !atomic.CompareAndSwapInt32(&o.state, running, scheduledAndRunning) {
			if !atomic.CompareAndSwapInt32(&o.state, notScheduled, running) {
				return
			}
		} else {
			return
		}
	}
	go func() {
		o.body()
		if !atomic.CompareAndSwapInt32(&o.state, running, notScheduled) {
			o.body()
			atomic.StoreInt32(&o.state, notScheduled)
		}
	}()
}
