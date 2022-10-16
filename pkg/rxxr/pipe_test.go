package rxxr

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"testing"
)

func TestPipe(t *testing.T) {
	if _, ok := any(new(pipe[any])).(Pipe[any]); !ok {
		t.Fail()
	}
}

func watchWithTimeout(t *testing.T, done <-chan struct{}) {
	t.Helper()
	ctx := context.Background()
	if deadline, ok := t.Deadline(); ok {
		var cancel func()
		ctx, cancel = context.WithDeadline(ctx, deadline)
		defer cancel()
	}

	select {
	case <-done:
	case <-ctx.Done():
		t.Fail()
	}
}

func Test_value_Publish(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		var (
			p           = New[int](nil)
			i           = rand.Int()
			ctx, cancel = context.WithCancel(context.Background())
		)

		p.Subscribe(func(v int) {
			defer cancel()
			if v != i {
				t.Logf("expected %d, got %d", i, v)
				t.Fail()
			}
		})
		p.Publish(i)
		watchWithTimeout(t, ctx.Done())
	})

	t.Run("publish many subscribers", func(t *testing.T) {
		p := New[int](nil)
		wg := new(sync.WaitGroup)
		count := math.MaxInt16
		ctx, cancel := context.WithCancel(context.Background())

		for i := 0; i <= count; i++ {
			wg.Add(1)
			p.Subscribe(func(v int) {
				defer wg.Done()
			})
		}

		go func() {
			t.Helper()
			defer cancel()
			wg.Wait()
		}()

		p.Publish(0)
		watchWithTimeout(t, ctx.Done())
	})
}

func Test_value_Subscribe(t *testing.T) {
	p := New[any](nil)

	t.Run("basic subscription", func(t *testing.T) {
		s := p.Subscribe(func(any) {})
		if s == nil {
			t.Fail()
		} else if !s.Subscribed() {
			t.Fail()
		}
	})

	t.Run("nil function subscription", func(t *testing.T) {
		s := p.Subscribe(nil)
		if s == nil {
			t.Fail()
		} else if s.Subscribed() {
			t.Fail()
		}
	})
}

func Test_value_Unsubscribe(t *testing.T) {
	t.Run("basic unsubscribe", func(t *testing.T) {
		p := New[any](nil)
		s := p.Subscribe(func(any) {})
		s.Unsubscribe()
		if s.Subscribed() {
			t.Fail()
		}
	})

	t.Run("avoid weird state", func(t *testing.T) {
		p := New[any](nil)
		rawP := p.(*pipe[any])
		subscriptions := &rawP.subscriptions
		subscriptions.values = append(subscriptions.values, nil)
		s := p.Subscribe(func(any) {
			t.Log("called subscription")
			t.Fail()
		})
		s.(*subscription[any]).subbed.Store(false)

		ctx, cancel := context.WithCancel(context.Background())
		s = p.Subscribe(func(any) {
			defer cancel()
			subscriptions.lock.Lock()
			defer subscriptions.lock.Unlock()
		})
		p.Publish(nil)
		watchWithTimeout(t, ctx.Done())

		s.Unsubscribe()
		if len(subscriptions.values) > 1 {
			t.Log("did not prune subscriptions list")
			t.Fail()
		}
	})
}

func Test_value_Value(t *testing.T) {
	t.Run("no value", func(t *testing.T) {
		p := New[any](nil)
		if _, ok := p.Value(); ok {
			t.Fail()
		}
	})

	t.Run("has value", func(t *testing.T) {
		i := rand.Int()
		p := New[int](new(PipeConfig[int]).SetInitialValue(i))
		if v, ok := p.Value(); !ok {
			t.Fail()
		} else if v != i {
			t.Fail()
		}
	})
}
