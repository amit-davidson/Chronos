package main

import "sync"

var a int

func main() {
	var mu sync.Mutex
	mu.Lock()
	a = 5
	_ = a
}
