# rxxr

Allows values to be published and subscribed to.
Subscribers will receive the most recently published value, if one exists.

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
// When a function is subscribed, 
// call it with the current value if one exists.
cfg.SendOnSubscribe(true)

// To use the config:
pipe := rxxr.New[int](cfg)
// or
pipe := cfg.Build()
```

### Subscribe

```go
pipe := rxxr.New[int]()

sub := pipe.Subscribe(func(i int) {
	// Do something
})

// Do something

pipe.Unsubscribe(sub)
```

### Publish

```go
pipe := rxxr.New[int]()
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
pipe := rxxr.New[int]()
fmt.Println(pipe.Value())

pipe.Publish(42)
fmt.Println(pipe.Value())

// 0 false
// 4 true
```

### Close

```go
pipe := rxxr.New[int]()

// Do something ...

pipe.Close()
```