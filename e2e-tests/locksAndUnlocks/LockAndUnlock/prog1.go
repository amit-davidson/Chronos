package pkg

import (
	"sync"
)

var a int

func main() {
	mutex := sync.Mutex{}
	mutex.Lock()
	mutex.Unlock()
	a = 5
}
