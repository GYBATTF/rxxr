// Package rxxr is a set of interfaces defining a ways to publish and subscribe to values.
// It also contains an implementation of a simple pipe that takes values and publishes them to multiple subscribers.
package rxxr

import (
	"github.com/google/uuid"
)

type (
	// Subscription is used to identify a unique subscription
	Subscription interface {
		GetID() uuid.UUID
		Subscribed() bool
	}

	Subscribable[T any] interface {
		// Subscribe registers a function to receive values from this pipe
		Subscribe(func(T)) Subscription
		// Unsubscribe unregisters a function from receiving values
		Unsubscribe(Subscription)
	}

	Publishable[T any] interface {
		// Publish publishes values to a pipe
		Publish(v ...T)
	}

	// Pipe is a type that can publish values to subscribers
	Pipe[T any] interface {
		Subscribable[T]
		Publishable[T]
		// Close closes the pipe and unsubscribes all subscribers
		Close()
		// Value returns the currently stored value, and if one exists
		Value() (t T, ok bool)
	}
)
