package mycontext

import (
	"errors"
	"reflect"
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
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
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
		ctx.cancel(errCancelled)
	}

	// Element cancel propagation
	go func() {
		select {
		case <-parent.Done():
			if ctx.err != nil {
				ctx.cancel(parent.Err())
			}
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}

func (ctx *cancelCtx) cancel(err error) {
	if ctx.err != nil {
		return
	}
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.err = err
	close(ctx.done)
}

type deadlineCtx struct {
	*cancelCtx
	deadline time.Time
}

func (ctx *deadlineCtx) Deadline() (time.Time, bool) {
	return ctx.deadline, true
}

var errDeadline = errors.New("deadline exceeded")

func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	cctx, cancel := WithCancel(parent)

	ctx := &deadlineCtx{
		cancelCtx: cctx.(*cancelCtx),
		deadline:  d,
	}

	timer := time.AfterFunc(time.Until(ctx.deadline), func() {
		ctx.cancel(errDeadline)
	})

	return ctx, func() {
		timer.Stop()
		cancel()
	}
}

var errTimeout = errors.New("timelimit exceeded")

type timeoutCtx struct {
	*deadlineCtx
}

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	dctx, _ := WithDeadline(parent, time.Now().Add(timeout))

	ctx := &timeoutCtx{
		deadlineCtx: dctx.(*deadlineCtx),
	}

	return ctx, func() { ctx.cancel(errTimeout) }
}

type valueCtx struct {
	*cancelCtx
	key   interface{}
	value interface{}
}

func (ctx *valueCtx) Value(key interface{}) interface{} {
	if ctx.key == key {
		return ctx.value
	}
	return nil
}

func WithValue(parent Context, key interface{}, value interface{}) (Context, CancelFunc) {
	if key == nil || reflect.TypeOf(key).Comparable() {
		panic("Key cannot be null or non-comparable type.")
	}

	ctx, cancel := WithCancel(parent)

	vctx := &valueCtx{
		cancelCtx: ctx.(*cancelCtx),
		key:       key,
		value:     value,
	}

	return vctx, cancel
}
