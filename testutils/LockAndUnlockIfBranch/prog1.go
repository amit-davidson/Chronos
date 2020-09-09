package pkg

import (
	"math/rand"
	"sync"
)

var a int
var cond = false
func main() {
	mutex := sync.Mutex{}
	if rand.Int() > 0 {
		mutex.Lock()
	}
	a = 5
	if rand.Int() > 0 {
		mutex.Unlock()
	}
	a = 6
}
