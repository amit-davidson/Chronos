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
		a = 6
	}()
	defer func() {
		defer mutex.Lock()
		defer mutex.Unlock()
	}()
}
