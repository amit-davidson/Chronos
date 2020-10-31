package main

import (
	"runtime"
	"time"
)

func main() {
	// It is possible that a future implementation
	// of memory allocation will ruin this test.
	// Increasing n might help in this case, so
	// this test is a bit more generic than most of the
	// others.
	const n = 2
	done := make(chan bool, n)
	empty := func(p *int) {}
	for i := 0; i < n; i++ {
		ms := i
		go func() {
			<-time.After(time.Duration(ms) * time.Millisecond)
			runtime.GC()
			var x int
			empty(&x) // x goes to the heap
			done <- true
		}()
	}
	for i := 0; i < n; i++ {
		<-done
	}
}
