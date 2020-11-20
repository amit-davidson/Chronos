package main

import "sync"

var mu sync.Mutex

func main() {
	f()
}

func f() {
	g()
	mu.Unlock()
}

func g() {
	mu.Lock()
}
