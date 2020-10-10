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
		defer mutex.Lock()
	}
	defer func() {
		a = 5
	}()
	if rand.Int() > 0 {
		defer mutex.Unlock()
	}
	a = 6
}
