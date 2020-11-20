package pkg

import (
	"math/rand"
	"sync"
)

var a = make(map[string]interface{}, 0)

func main() {
	mutex := sync.Mutex{}
	if rand.Int() > 0 {
		if rand.Int() > 0 {
			mutex.Lock()
		} else {
			mutex.Lock()
		}
	} else {
		mutex.Lock()
	}
	a["1"] = make(map[string]interface{}, 0)
}
