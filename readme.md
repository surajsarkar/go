# Project

* Define a interface
```go
type y interface {}; // where y in interface name
```
* Define `struct`
```go
type a struct {} // a is struct name
```
* Adding methods to struct
```go
func (x structname) funcname (<input> <inputtype>) <output_type> {
    // code to execute
}
```

* condition 
```go 
if (condition) {
    // code to execute if condition is true
} else if (another_condition){
    // execute if another_condition true
}else {
    // default case
}
```

* creating map 
```go
map_name := map[<key_dtype>]<val_dtype>{values...}
//eg
score := map[string]int{"Suraj": 100, "Opponent": 1}
```
## Channels
### Things to remember
* Sending on a closed channel will cause a panic
> when we use channel what ever comes in channel will not be sequential.


## Mutex
* helps in locking machanism so if you are accessing a value across a routines, to remove conflict we use mutex so a single routine can use the channel at a time. 
* to lock use : `sync.Mutex.Lock()`
* to lock use : `sync.Mutex.Unlock()`


## Contex
* Context type, which carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes
* when we create a `context` it returns a context object and **cancel** func, which when called sends a message to `context.Done()` channel. which tells if the channel was closed.
## Things to cover with projects
* context


some important points 
```go 
//assigning pointer of a var to another var
i := 5
p := &i

// here p gets the value of i's memory address
```

To create a slice dinamically sized array 
```go 
slice := make([]int, 7) // dtype, length, capacity (optional)
//append to slice 
append(slice, element_to_append)
```
> Then we will start building project