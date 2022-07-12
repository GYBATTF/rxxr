package rxxr

import (
	"github.com/google/uuid"
	"sync"
)

type (
	pipe[T any] struct {
		cfg Config[T]

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

// New creates a simple pipe that transmits received values to all subscribers
func New[T any](cfg *Config[T]) Pipe[T] {
	if cfg == nil {
		cfg = &Config[T]{}
	}

	v := &pipe[T]{
		cfg:           *cfg,
		valueChange:   make(chan T),
		subscriptions: map[uuid.UUID]*subscription[T]{},
	}

	if cfg.initialValueSet {
		v.value = cfg.initialValue
		v.isSet = true
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
	if v, ok := val.Value(); val.cfg.sendOnSubscribe && ok {
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

// Config is used to configure a pipe.
// Can either be used as a parameter to New or on its own as a builder.
type Config[T any] struct {
	initialValue    T
	initialValueSet bool

	sendOnSubscribe bool
}

// Build creates a new pipe with this config.
func (c *Config[T]) Build() Pipe[T] {
	return New[T](c)
}

// SetInitialValue Initializes a pipe with the provided value.
func (c *Config[T]) SetInitialValue(v T) *Config[T] {
	c.initialValueSet = true
	c.initialValue = v
	return c
}

// SendOnSubscribe sets if the current value should be sent to new subscribers,
// if a value is not sent the subscriber will receive no values.
func (c *Config[T]) SendOnSubscribe(b bool) *Config[T] {
	c.sendOnSubscribe = b
	return c
}
