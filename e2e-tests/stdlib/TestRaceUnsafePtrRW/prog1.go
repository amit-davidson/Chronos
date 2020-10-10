package main

import "unsafe"

func main() {
	var x, y, z int
	x, y, z = 1, 2, 3
	var p unsafe.Pointer = unsafe.Pointer(&x)
	ch := make(chan bool, 1)
	go func() {
		p = (unsafe.Pointer)(&z)
		ch <- true
	}()
	y = *(*int)(p)
	x = y
	<-ch
}
