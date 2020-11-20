package Lock

import (
	"sync"
)

var a int

type obj3 struct {
	obj2 obj2
}
type obj2 struct {
	obj1 obj1
}

type obj1 struct {
	a     int
	mutex sync.Mutex
}

func main() {
	obj3 := obj3{obj2: obj2{obj1: obj1{1, sync.Mutex{}}}}
	obj3.obj2.obj1.mutex.Lock()
	obj3.obj2.obj1.mutex.Unlock()
	a = 5
}
