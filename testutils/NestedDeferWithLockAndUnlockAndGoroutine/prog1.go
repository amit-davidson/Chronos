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
		mutex.Lock()
		a = 6
		defer mutex.Unlock()
	}()
	mutex.Lock()
	defer func() {
		go func() {
			a=7
			defer mutex.Unlock()
		}()
		//defer mutex.Unlock()
	}()
}
