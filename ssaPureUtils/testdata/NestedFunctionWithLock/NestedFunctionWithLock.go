package functionWithLock

import "sync"

var a int

func main() {
	f()
}

func f() {
	var mu sync.Mutex
	mu.Lock()
	a = 5
	_ = a
}