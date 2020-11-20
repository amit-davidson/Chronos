package main

import "sync"

func main() {
	ch := make(chan bool, 1)
	var mu sync.Mutex
	x := 0
	_ = x
	go func() {
		mu.Lock()
		x = 42
		mu.Unlock()
		ch <- true
	}()
	x = func(mu *sync.Mutex) int {
		mu.Lock()
		return 43
	}(&mu)
	mu.Unlock()
	<-ch
}
