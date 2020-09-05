package Lock

import (
	"sync"
)

var a int

func fn1() {
	mutex := sync.Mutex{}
	mutex.Lock()
	a = 5
}
