package rxxr

import (
	"math"
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
	p := New[any]()
	if p == nil {
		t.Fail()
	} else {
		p.Close()
	}
}

func TestWith(t *testing.T) {
	i := math.Max(rand.Float64(), 1)
	p := With[float64](i)
	if p == nil {
		t.FailNow()
	}
	defer p.Close()

	if v, _ := p.Value(); v != i {
		t.Fail()
	}
}

func Test_value_Close(t *testing.T) {
	p := New[any]()
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
	case <-time.After(5 * time.Second):
	}
}

func Test_value_Publish(t *testing.T) {
	p := New[int]()

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
	p := New[any]()
	defer p.Close()
	s := p.Subscribe(func(any) {})
	if !s.Subscribed() {
		t.Fail()
	}
}

func Test_value_Unsubscribe(t *testing.T) {
	p := New[any]()
	defer p.Close()
	s := p.Subscribe(func(any) {})
	p.Unsubscribe(s)
	if s.Subscribed() {
		t.Fail()
	}
}

func Test_value_Value(t *testing.T) {
	t.Run("no value", func(t *testing.T) {
		p := New[any]()
		defer p.Close()
		if _, ok := p.Value(); ok {
			t.Fail()
		}
	})

	t.Run("has value", func(t *testing.T) {
		i := rand.Int()
		p := With[int](i)
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
