package main

import (
	"math/rand"
	"sync"
)

var a int

func main() {
	var mu sync.Mutex
	if rand.Int() > 0 {
		mu.Unlock()
	} else {
		mu.Lock()
	}
	a = 5
	_ = a
}
