package pkg

import (
	"math/rand"
	"sync"
)

var a = make(map[string]interface{}, 0)
func fn1() {
	mutex := sync.Mutex{}
	if rand.Int() > 0 {
		mutex.Lock()
	}
	a["1"] = make(map[string]interface{}, 0)
	val := a["1"]
	val.(map[string]interface{})["2"] = "3"
	if rand.Int() > 0 {
		mutex.Unlock()
	}
}
