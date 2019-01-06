# dispatcher

## about
dispatcher is a function scheduler that dispatches a limited number of functions and keeps rest
of the functions in a list for dispatching later.

## use
import pkg
```go
import github.com/sdeoras/dispatcher
```

dispatch function
```go
d := dispatcher.New(5) // dispatch a max of 5 functions concurrently
d.Do(func(){
	fmt.Println("dispatched")
})
```

wait for completion
```go
for d.IsRunning() {
	time.Sleep(time.Second)
}
```