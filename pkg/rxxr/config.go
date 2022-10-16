package rxxr

// PipeConfig is used to configure a pipe.
// Can either be used as a parameter to New or on its own as a builder.
type PipeConfig[T any] struct {
	// InitialValue Initializes a pipe with the provided value.
	InitialValue *T

	// SendOnSubscribe sets if the current value should be sent to new subscribers,
	// if a value is not sent the subscriber will receive no values.
	SendOnSubscribe bool

	PanicLogger func(r any)
}

func (c *PipeConfig[T]) SetInitialValue(v T) *PipeConfig[T] {
	c.InitialValue = &v
	return c
}

func (c *PipeConfig[T]) SetSendOnSubscribe(b bool) *PipeConfig[T] {
	c.SendOnSubscribe = b
	return c
}

func (c *PipeConfig[T]) SetPanicLogger(fn func(r any)) *PipeConfig[T] {
	c.PanicLogger = fn
	return c
}
