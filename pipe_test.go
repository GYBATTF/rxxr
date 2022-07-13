package rxxr

import (
	"math/rand"
	"testing"
	"time"
)

func TestPipe(t *testing.T) {
	if _, ok := any(new(pipe[any])).(Pipe[any]); !ok {
		t.Fail()
	}
}

func TestNew(t *testing.T) {
	t.Run("new with nil config", func(t *testing.T) {
		p := New[any](nil)
		if p == nil {
			t.Fail()
		} else {
			p.Close()
		}
	})

	t.Run("new with config", func(t *testing.T) {
		p := New[any](&Config[any]{})
		if p == nil {
			t.Fail()
		} else {
			p.Close()
		}
	})

	t.Run("config build method", func(t *testing.T) {
		p := (&Config[any]{}).Build()
		if p == nil {
			t.Fail()
		} else {
			p.Close()
		}
	})
}

func Test_value_Close(t *testing.T) {
	passDelay := 5 * time.Second

	p := New[any](nil)
	c := make(chan int)
	defer close(c)
	p.Subscribe(func(any) {
		c <- 0
	})
	p.Close()
	p.Publish(1)

	select {
	case <-c:
		t.Fail()
	case <-time.After(passDelay):
	}
}

func Test_value_Publish(t *testing.T) {
	p := New[int](nil)

	i := rand.Int()
	received := make(chan bool)
	defer close(received)

	p.Subscribe(func(v int) {
		if v != i {
			t.Logf("expected %d, got %d", i, v)
			t.Fail()
		}
		received <- true
	})

	p.Publish(i)

	select {
	case <-received:
	case <-time.After(5 * time.Second):
		t.Log("time out")
		t.Fail()
	}
	p.Close()
}

func Test_value_Subscribe(t *testing.T) {
	p := New[any](nil)
	defer p.Close()
	s := p.Subscribe(func(any) {})
	if s == nil {
		t.Fail()
	} else if !s.Subscribed() {
		t.Fail()
	}
}

func Test_value_Unsubscribe(t *testing.T) {
	p := New[any](nil)
	defer p.Close()
	s := p.Subscribe(func(any) {})
	p.Unsubscribe(s)
	if s.Subscribed() {
		t.Fail()
	}
}

func Test_value_Value(t *testing.T) {
	t.Run("no value", func(t *testing.T) {
		p := New[any](nil)
		defer p.Close()
		if _, ok := p.Value(); ok {
			t.Fail()
		}
	})

	t.Run("has value", func(t *testing.T) {
		i := rand.Int()
		p := New[int](new(Config[int]).SetInitialValue(i))
		defer p.Close()
		if v, ok := p.Value(); !ok {
			t.Fail()
		} else if v != i {
			t.Fail()
		}
	})
}

func Test_subscription(t *testing.T) {
	if _, ok := any(new(subscription[any])).(Subscription); !ok {
		t.Fail()
	}
}
