package main

import (
	"fmt"
	"sync"
)
import "time"

type myType struct {
	A int
}

var mutex = sync.Mutex{}
func main() {
	x := new(myType)
	c := make(chan int, 100)
	go func() {
		for {
			mutex.Lock()
			x = new(myType) // write to x
			c <- 0
			<-c
		}
	}()
	for i := 0; i < 4; i++ {
		go func() {
			for {
				_ = *x // if exists a race condition, `*x` will visit a wrong memory address, and will panic
				c <- 0
				<-c
			}
		}()
	}
	time.Sleep(time.Second * 10)
	fmt.Println("end")
}
