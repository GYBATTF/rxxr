package rxxr

import "testing"

func Test_subscription(t *testing.T) {
	if _, ok := any(new(subscription[any])).(Subscription); !ok {
		t.Fail()
	}
}

func TestSubscription_call(t *testing.T) {
	p := New[any](nil)
	sub := newSubscription[any](func(any) {}, p.(*pipe[any]))
	sub.fn = nil
	if sub.call(nil) {
		t.Log("called nil fn")
		t.Fail()
	}
	sub.fn = func(any) {}
	sub.subbed.Store(false)
	if sub.call(nil) {
		t.Log("called unsubbed")
		t.Fail()
	}
}
