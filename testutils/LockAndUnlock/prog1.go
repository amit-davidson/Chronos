package pkg

import (
	"sync"
)

var a int

func fn1() {
	mutex := sync.Mutex{}
	mutex.Lock()
	mutex.Unlock()
	a = 5
}
