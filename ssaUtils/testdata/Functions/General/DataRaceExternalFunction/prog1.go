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
	fn1()
}

func fn1() {
	if a > 4 {
		goto End
	}
	return
End:
	a = 6
}
