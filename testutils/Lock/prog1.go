package Lock

import (
	"sync"
)

var a int

func main() {
	mutex := sync.Mutex{}
	mutex.Lock()
	a = 5
}
