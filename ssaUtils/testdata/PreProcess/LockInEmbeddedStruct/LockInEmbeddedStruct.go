package functionWithLock

import "sync"

var a int

type B struct {
	C
}

type C struct {
	Lock sync.Mutex
}

func main() {
	b := B{C: C{
		sync.Mutex{},
	}}
	b.Lock.Lock()
	a = 5
	_ = a
}
