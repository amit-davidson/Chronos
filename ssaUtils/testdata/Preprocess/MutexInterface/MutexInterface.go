package main

import "sync"

// A Locker represents an object that can be locked and unlocked.
type Locker interface {
	Lock()
	Unlock()
}

type A struct {
	Lock Locker
}

func main() {
	b := A{Lock: &sync.Mutex{}}
	b.Lock.Lock()
}
