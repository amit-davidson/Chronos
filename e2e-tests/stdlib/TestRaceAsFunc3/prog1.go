package main

import "sync"

func main() {
	c := make(chan bool, 1)
	var mu sync.Mutex
	x := 0
	go func() {
		func(x int) {
			mu.Lock()
		}(x) // Read of x must be outside of the mutex.
		mu.Unlock()
		c <- true
	}()
	mu.Lock()
	x = 42
	mu.Unlock()
	<-c
}
