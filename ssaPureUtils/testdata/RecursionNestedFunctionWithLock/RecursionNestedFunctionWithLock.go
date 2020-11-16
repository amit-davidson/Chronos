package functionWithLock

import "sync"

var mu sync.Mutex

func main() {
	f()
}

func f() {
	mu.Lock()
	f()
}