package main

import "fmt"

// emptyFunc should not be inlined.
func emptyFunc(x int) {
	if false {
		fmt.Println(x)
	}
}

func main() {
	var x int
	ch := make(chan bool, 1)
	go func() {
		emptyFunc(x)
		ch <- true
	}()
	x = 1
	<-ch
}
