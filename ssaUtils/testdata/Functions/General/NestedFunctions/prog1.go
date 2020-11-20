package main

import (
	"sync"
)

var a int
var mutex1 = sync.Mutex{}

func main() {
	mutex2 := sync.Mutex{}
	mutex1.Lock()
	mutex2.Lock()
	a = 5
	fn2()
}

func fn2() {
	mutex2 := sync.Mutex{}
	mutex2.Lock()
	a = 6
	mutex2.Unlock()
	mutex1.Unlock()
}
