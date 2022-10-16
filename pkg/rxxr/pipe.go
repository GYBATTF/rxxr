// Package rxxr is a set of interfaces defining a ways to publish and subscribe to values.
// It also contains an implementation of a simple pipe that takes values and publishes them to multiple subscribers.
package rxxr

import (
	"sync"
)

type (
	Subscribable[T any] interface {
		// Subscribe registers a function to receive values from this pipe
		Subscribe(func(T)) Subscription
	}

	Publishable[T any] interface {
		// Publish publishes values to a pipe
		Publish(v T)
	}

	// Pipe is a type that can publish values to subscribers
	Pipe[T any] interface {
		Subscribable[T]
		Publishable[T]
		// Value returns the currently stored value, and if one exists
		Value() (t T, ok bool)
	}
)

type pipe[T any] struct {
	cfg PipeConfig[T]

	value struct {
		value T
		isSet bool
		lock  sync.RWMutex
	}

	subscriptions struct {
		values []*subscription[T]
		lock   sync.Mutex
	}
}

// New creates a simple pipe that transmits received values to all subscribers
func New[T any](cfg *PipeConfig[T]) Pipe[T] {
	p := new(pipe[T])
	if cfg != nil {
		p.cfg = *cfg
	}

	if p.cfg.InitialValue != nil {
		p.setValue(*p.cfg.InitialValue)
	}
	return p
}

func (val *pipe[T]) Publish(v T) {
	val.setValue(v)
	go val.publish(v)
}

func (val *pipe[T]) publish(v T) {
	val.subscriptions.lock.Lock()
	defer val.subscriptions.lock.Unlock()
	for i := 0; i < len(val.subscriptions.values); {
		if e := val.getKey(i); e != nil {
			go e.call(v)
			i++
		}
	}
}

func (val *pipe[T]) setValue(v T) {
	val.value.lock.Lock()
	defer val.value.lock.Unlock()
	val.value.value, val.value.isSet = v, true
}

func (val *pipe[T]) getKey(i int) *subscription[T] {
	e := val.subscriptions.values[i]
	if e == nil || !e.Subscribed() {
		val.subscriptions.values = append(val.subscriptions.values[:i], val.subscriptions.values[i+1:]...)
		return nil
	}
	return e
}

func (val *pipe[T]) Subscribe(fn func(T)) Subscription {
	return val.checkSubscription(newSubscription[T](fn, val))
}

func (val *pipe[T]) checkSubscription(s *subscription[T]) Subscription {
	if s.fn == nil {
		s.Unsubscribe()
	} else {
		val.addSubscription(s)
		go val.sendInitial(func(t T) { s.call(t) })
	}
	return s
}

func (val *pipe[T]) addSubscription(s *subscription[T]) {
	val.subscriptions.lock.Lock()
	defer val.subscriptions.lock.Unlock()
	val.subscriptions.values = append(val.subscriptions.values, s)
}

func (val *pipe[T]) sendInitial(fn func(T)) {
	if v, ok := val.Value(); val.cfg.SendOnSubscribe && ok {
		go fn(v)
	}
}

func (val *pipe[T]) Value() (t T, ok bool) {
	val.value.lock.RLock()
	defer val.value.lock.RUnlock()
	return val.value.value, val.value.isSet
}
