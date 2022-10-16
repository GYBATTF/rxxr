# rxxr

Allows values to be published and subscribed to.
Subscribers will receive the most recently published value, if one exists.

Subscribed functions are run in their own goroutine.
Calling methods on the pipe are thread safe, but published values will not be locked.

### Constructor

```go
// Provide nil as config for defaults
pipe := rxxr.New[int](nil)
```

### Configure

```go
cfg := &rxxr.Config[int]{}

// Create a pipe already containing a value
cfg.SetInitialValue(42)
// -- or --
i := 0
cfg.InitialValue = &i

// When a function is subscribed, 
// call it with the current value if one exists.
cfg.SendOnSubscribe = true

// To use the config:
pipe := rxxr.New[int](cfg)
```

### Subscribe

```go
pipe := rxxr.New[int](nil)

sub := pipe.Subscribe(func(i int) {
	// Do something
})

// Do something

sub.Unsubscribe()
```

### Publish

```go
pipe := rxxr.New[int](nil)
pipe.Subscribe(func (i int) {
	fmt.Printf("received: \n", i)
})
pipe.Publish(1)
pipe.Publish(2, 3)
pipe.Publish(5, 8, 13)

// received: 1
// received: 2
// received: 3
// received: 5
// received: 8
// received: 13
```

### Get Value

```go
pipe := rxxr.New[int](nil)
fmt.Println(pipe.Value())

pipe.Publish(42)
fmt.Println(pipe.Value())

// 0 false
// 4 true
```