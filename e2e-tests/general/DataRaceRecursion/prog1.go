package main

import (
	"sync"
)

var a int
var mutex1 = sync.Mutex{}

func main() {
	go func() {
		a = 5
	}()
	fn1(0)
}

func fn1(counter int) {
	if counter >= 10 {
		return
	}
	if counter == 8 {
		a = 6
	}
	fn1(counter)
}
