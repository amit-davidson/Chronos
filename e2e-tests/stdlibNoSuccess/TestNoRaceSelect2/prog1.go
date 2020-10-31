package main

import "runtime"

func main() {
	var x int
	_ = x
	compl := make(chan bool)
	c := make(chan bool)
	c1 := make(chan bool)
	go func() {
		select {
		case <-c:
		case <-c1:
		}
		x = 1
		compl <- true
	}()
	x = 2
	close(c)
	runtime.Gosched()
	<-compl
}
