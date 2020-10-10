package main

import "sync"

func main() {
	c := make(chan bool, 1)
	var mu sync.Mutex
	x := 0
	_ = x
	go func() {
		x = func() int { // Write of x must be under the mutex.
			mu.Lock()
			return 42
		}()
		mu.Unlock()
		c <- true
	}()
	mu.Lock()
	x = 42
	mu.Unlock()
	<-c
}
