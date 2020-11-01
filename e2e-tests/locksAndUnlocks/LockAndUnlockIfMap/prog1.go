package pkg

import (
	"math/rand"
	"sync"
)

var a = make(map[string]interface{}, 0)

func main() {
	mutex := sync.Mutex{}
	if rand.Int() > 0 {
		mutex.Lock()
	}
	a["1"] = make(map[string]interface{}, 0)
	val := a["1"]
	val.(map[string]interface{})["2"] = "2"
	if rand.Int() > 0 {
		mutex.Unlock()
	}
	a["2"] = make(map[string]interface{}, 0)
}
