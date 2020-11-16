package functionWithLock

import (
	"math/rand"
	"sync"
)

var a int

func main() {
	var mu sync.Mutex
	mu.Lock()
	if rand.Int() > 0 {
		goto Unlock
	}
	goto AfterUnlock
Unlock:
	mu.Unlock()
AfterUnlock:
	a = 5
	_ = a
}
