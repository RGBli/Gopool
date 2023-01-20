A simple goroutine pool
---

### Usage
lazy load pool
```go
p := NewPoolLazyWorker(1024)
p.Schedule(func() {
    fmt.Println("hello go pool")
})
time.Sleep(time.Second * 2)
fmt.Println("done")
```

hungry load pool
```go
p := NewPoolLazyWorker(1024, 512)
p.Schedule(func() {
    fmt.Println("hello go pool")
})
time.Sleep(time.Second * 2)
fmt.Println("done")
```

pool with timeout
```go
p := NewPoolLazyWorker(1024)
p.Schedule(func() {
    fmt.Println("hello go pool")
}, time.Second * 1)
time.Sleep(time.Second * 2)
fmt.Println("done")
```