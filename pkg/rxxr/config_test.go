package rxxr

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("new with nil config", func(t *testing.T) {
		p := New[any](nil)
		if p == nil {
			t.Fail()
		}
	})

	t.Run("new with config", func(t *testing.T) {
		p := New[any](&PipeConfig[any]{})
		if p == nil {
			t.Fail()
		}
	})

	t.Run("new with initial value", func(t *testing.T) {
		p := New[any](new(PipeConfig[any]).SetInitialValue(5))
		v, ok := p.Value()
		if p == nil || v != 5 || !ok {
			t.Fail()
		}
	})

	t.Run("new with send on subscribe", func(t *testing.T) {
		p := New[int](new(PipeConfig[int]).SetSendOnSubscribe(true).SetInitialValue(5))
		if p == nil {
			t.Fail()
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		p.Subscribe(func(i int) {
			defer cancel()
			if i != 5 {
				t.Fail()
			}
		})
		watchWithTimeout(t, ctx.Done())
	})

	t.Run("new with panic handler", func(t *testing.T) {
		var (
			ctx, cancel = context.WithCancel(context.Background())
			p           = New[int](new(PipeConfig[int]).SetSendOnSubscribe(true).SetInitialValue(5).SetPanicLogger(func(r any) {
				defer cancel()
				if i, ok := r.(int); !ok || i != 5 {
					t.Fail()
				}
			}))
		)
		defer cancel()

		if p == nil {
			t.Fail()
			return
		}
		p.Subscribe(func(i int) { panic(i) })
		watchWithTimeout(t, ctx.Done())
	})
}
