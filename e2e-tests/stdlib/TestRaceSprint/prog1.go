package main

import "fmt"

func main() {
	var x int
	ch := make(chan bool, 1)
	go func() {
		fmt.Sprint(x)
		ch <- true
	}()
	x = 1
	<-ch
}
