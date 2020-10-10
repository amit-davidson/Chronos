package main

import "sync"

type Point struct {
	x, y int
}
var mutex sync.Mutex

func main() {
	c := make(chan bool)
	var x interface{}
	go func() {
		x = &Point{}
		c <- true
	}()
	mutex.Lock()
	_ = x
	<-c
}
