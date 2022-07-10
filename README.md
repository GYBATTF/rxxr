# rxxr

Allows values to be published and subscribed to.
Subscribers will receive the most recently published value, if one exists.

### Constructors

```go
pipe := rxxr.New[int]()
pipe.Subscribe(func (i int) {
	// Will not run until a value is published
})
```

```go
pipe := rxxr.With[int](42)
pipe.Subscribe(func (i int) {
    // Will immediately run with initial value
	fmt.Println(42)
})

// 42
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