package rxxr

import (
	"github.com/google/uuid"
	"sync"
)

type (
	pipe[T any] struct {
		value       T
		isSet       bool
		valueChange chan T

		closed        bool
		subscriptions map[uuid.UUID]*subscription[T]

		locks struct {
			closed        sync.RWMutex
			value         sync.RWMutex
			subscriptions sync.RWMutex
		}
	}
)

// With creates a new pipe already containing a value
func With[T any](initial T) Pipe[T] {
	val := New[T]().(*pipe[T])
	val.setValue(initial)
	return val
}

// New creates a simple pipe that transmits received values to all subscribers
func New[T any]() Pipe[T] {
	v := &pipe[T]{
		valueChange:   make(chan T),
		subscriptions: map[uuid.UUID]*subscription[T]{},
	}
	go v.watchForValueChange()
	return v
}

func (val *pipe[T]) watchForValueChange() {
	for v := range val.valueChange {
		val.sendNewValue(v)
	}
}

func (val *pipe[T]) sendNewValue(v T) {
	val.locks.closed.RLock()
	defer val.locks.closed.RUnlock()
	if val.closed {
		return
	}
	val.setValue(v)
	val.sendValueToAllSubscribers(v)
}

func (val *pipe[T]) sendValueToAllSubscribers(v T) {
	val.locks.subscriptions.RLock()
	defer val.locks.subscriptions.RUnlock()
	for _, s := range val.subscriptions {
		s.pipe <- v
	}
}

func (val *pipe[T]) setValue(v T) {
	val.locks.value.Lock()
	defer val.locks.value.Unlock()
	val.value, val.isSet = v, true
}

func (val *pipe[T]) Publish(values ...T) {
	val.locks.closed.RLock()
	defer val.locks.closed.RUnlock()
	if !val.closed {
		for _, v := range values {
			val.valueChange <- v
		}
	}
}

func (val *pipe[T]) Subscribe(fn func(T)) Subscription {
	val.locks.closed.RLock()
	defer val.locks.closed.RUnlock()
	if val.closed || fn == nil {
		return new(subscription[T])
	}

	s := newSubscription[T](fn)
	if v, ok := val.Value(); ok {
		s.pipe <- v
	}

	val.locks.subscriptions.Lock()
	defer val.locks.subscriptions.Unlock()
	val.subscriptions[s.id] = s
	return s
}

func (val *pipe[T]) Unsubscribe(s Subscription) {
	val.locks.closed.RLock()
	defer val.locks.closed.RUnlock()
	if val.closed {
		return
	}

	if sub, ok := s.(*subscription[T]); ok && sub.subscribed {
		val.locks.subscriptions.Lock()
		defer val.locks.subscriptions.Unlock()
		delete(val.subscriptions, sub.id)

		sub.subscribed = false
		close(sub.pipe)
		sub.fn = nil
	}
}

func (val *pipe[T]) Close() {
	for _, v := range val.subscriptions {
		val.Unsubscribe(v)
	}

	val.locks.closed.Lock()
	defer val.locks.closed.Unlock()
	val.closed = true
	close(val.valueChange)
}

func (val *pipe[T]) Value() (t T, ok bool) {
	val.locks.value.RLock()
	defer val.locks.value.RUnlock()
	return val.value, val.isSet
}

type (
	subscription[T any] struct {
		id         uuid.UUID
		pipe       chan T
		subscribed bool
		fn         func(T)
	}
)

func newSubscription[T any](fn func(T)) *subscription[T] {
	return &subscription[T]{
		id:         uuid.New(),
		pipe:       listen[T](fn),
		fn:         fn,
		subscribed: true,
	}
}

func listen[T any](fn func(T)) (c chan T) {
	c = make(chan T)
	go func() {
		for v := range c {
			fn(v)
		}
	}()
	return
}

func (s *subscription[T]) GetID() uuid.UUID {
	return s.id
}

func (s *subscription[T]) Subscribed() bool {
	return s.subscribed
}
