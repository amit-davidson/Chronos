package main

type Point struct {
	x, y int
}

func main() {
	ch := make(chan byte, 1)
	var x byte
	go func(y byte) {
		_ = y
		ch <- 0
	}(x)
	x = 1
	<-ch
}
