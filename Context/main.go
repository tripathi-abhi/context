package context

import (
	"errors"
	"sync"
	"time"
)

// Creating context interface

type Context interface {
	Deadline() (time.Time, bool)
	Done() <-chan struct{}
	Err() error
	Value(interface{}) interface{}
}

type emptyCtx int64

type CancelFunc func()

func (emptyCtx) Deadline() (deadline time.Time, ok bool) { return }
func (emptyCtx) Done() <-chan struct{}                   { return nil }
func (emptyCtx) Err() error                              { return nil }
func (emptyCtx) Value(interface{}) interface{}           { return nil }

var (
	background = new(emptyCtx)
	todo       = new(emptyCtx)
)

func Background() Context {
	return background
}
func TODO() Context {
	return todo
}

type cancelCtx struct {
	Context
	done chan struct{}
	err  error
	mu   sync.Mutex
}

var errCancelled = errors.New("context cancelled")

func (ctx *cancelCtx) Err() error {
	return ctx.err
}

func (ctx *cancelCtx) Done() <-chan struct{} {
	return ctx.done
}

func WithCancel(parent Context) (Context, CancelFunc) {

	ctx := &cancelCtx{
		Context: parent,
		done:    make(chan struct{}),
	}

	cancel := func() {
		ctx.mu.Lock()
		defer ctx.mu.Unlock()
		ctx.err = errCancelled
		close(ctx.done)
	}

	return ctx, cancel
}
