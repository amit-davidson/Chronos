package main

import "sync"

var a int
type B struct {
	Lock sync.Mutex
}
func main() {
	b := B{Lock:sync.Mutex{}}
	b.Lock.Lock()
	a = 5
	_ = a
}
