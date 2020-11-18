package functionWithLock

import "sync"

var a int

func main() {
	var mu sync.Mutex
	mu.Unlock()
	a = 5
	_ = a
}
