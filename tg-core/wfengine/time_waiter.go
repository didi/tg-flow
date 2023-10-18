package wfengine

import (
"time"
)

type TimeWaiter struct {
	ch chan struct{}
	timeout int64
	deltaTimeout int64
}

func NewTimeWaiter(timeOut int64) *TimeWaiter {
	if timeOut <=0 {
		return nil
	}

	return &TimeWaiter{
		ch:      make(chan struct{}, 1),
		timeout: timeOut,
		deltaTimeout:0,
	}
}

func (t *TimeWaiter) Wait() bool {
	tm := time.NewTimer(time.Millisecond * time.Duration(t.timeout + t.deltaTimeout))
	select {
	case <- t.ch:
		tm.Stop()
		return true
	case <- tm.C:
		return false
	}
}

func (t *TimeWaiter) Done(){
	t.ch <- struct{}{}
}

func (t *TimeWaiter) AddTimeout(timeout int64){
	t.deltaTimeout = timeout
}