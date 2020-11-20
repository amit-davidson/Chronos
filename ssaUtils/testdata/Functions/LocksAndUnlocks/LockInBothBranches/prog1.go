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
	} else {
		mutex.Lock()
	}
	a = 5
}
