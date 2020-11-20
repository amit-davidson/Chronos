package main

import "sync"

var a int

func main() {
	var mu sync.Mutex
	mu.Lock()
	mu.Unlock()
	a = 5
	_ = a
}
