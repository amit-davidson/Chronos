package pkg

import (
	"sync"
)

var a = make(map[string]interface{}, 0)
func main() {
	mutex := sync.Mutex{}
	go func() {
		mutex.Lock()
		a["1"] = make(map[string]interface{}, 0)
	}()
	a["2"] = make(map[string]interface{}, 0)
}
