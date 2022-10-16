package rxxr

import "sync/atomic"

type (
	// Subscription is used to identify a unique subscription
	Subscription interface {
		Subscribed() bool
		// Unsubscribe unregisters a function from receiving values
		Unsubscribe()
	}

	subscription[T any] struct {
		subbed atomic.Bool
		fn     func(T)
		pipe   *pipe[T]
	}
)

func newSubscription[T any](fn func(T), p *pipe[T]) *subscription[T] {
	s := &subscription[T]{
		fn: fn,
	}
	s.subbed.Store(true)
	s.pipe = p
	return s
}

func (s *subscription[T]) Subscribed() bool {
	return s.subbed.Load()
}

func (s *subscription[T]) Unsubscribe() {
	if s.subbed.CompareAndSwap(true, false) {
		// Clear everything
		s.pipe = nil
		s.fn = nil
	}
}

func (s *subscription[T]) catch(b *bool) {
	if r := recover(); r != nil {
		p := s.pipe
		if p != nil && p.cfg.PanicLogger != nil {
			s.pipe.cfg.PanicLogger(r)
		}
		*b = false
	}
}

func (s *subscription[T]) call(v T) (b bool) {
	if fn := s.fn; s.Subscribed() && fn != nil {
		defer s.catch(&b)
		fn(v)
		return true
	}
	return false
}
