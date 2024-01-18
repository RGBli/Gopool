package main

import (
	"fmt"
	"time"
)

var ErrScheduleTimeout = fmt.Errorf("schedule err: timeout")

type Pool struct {
	semCh  chan struct{}
	workCh chan func()
}

func NewPoolLazyWorker(size int64) *Pool {
	p := &Pool{
		workCh: make(chan func()),
		semCh:  make(chan struct{}, size),
	}
	return p
}

func NewPoolHungryWorker(size, preload int64) *Pool {
	if preload > size {
		panic("preload goroutines more than total goroutines")
	}

	p := &Pool{
		workCh: make(chan func()),
		semCh:  make(chan struct{}, size),
	}

	var i int64
	for i = 0; i < preload; i++ {
		p.semCh <- struct{}{}
		go p.worker(func() {})
	}
	return p
}

func (p *Pool) Schedule(task func()) error {
	return p.schedule(task, nil)
}

func (p *Pool) ScheduleWithTimeout(task func(), timeout time.Duration) error {
	return p.schedule(task, time.After(timeout))
}

func (p *Pool) schedule(task func(), timeout <-chan time.Time) error {
	select {
	case <-timeout:
		return ErrScheduleTimeout
	case p.workCh <- task:
		fmt.Println("another task comes in work channel")
	case p.semCh <- struct{}{}:
		go p.worker(task)
	}
	return nil
}

func (p *Pool) worker(task func()) {
	defer func() {
		<-p.semCh
	}()

	for {
		task()
		task = <-p.workCh
	}
}
