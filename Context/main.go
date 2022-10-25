package context

import "time"

// Creating context interface

type Context interface {
	Deadline() (time.Time, bool)
	Done() <-chan struct{}
	Err() error
	Value(interface{}) interface{}
}

type emptyCtx int64

func Deadline() (time.Time, bool)   { return time.Time{}, false }
func Done() <-chan struct{}         { return nil }
func Err() error                    { return nil }
func Value(interface{}) interface{} { return nil }

var (
	background = new(emptyCtx)
	todo       = new(emptyCtx)
)

func Background() *emptyCtx {
	return background
}

func TODO() *emptyCtx {
	return todo
}
