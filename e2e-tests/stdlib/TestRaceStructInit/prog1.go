package main

func main() {
	type X struct {
		x, y int
	}
	c := make(chan bool, 1)
	y := 0
	go func() {
		y = 42
		c <- true
	}()
	x := X{x: y}
	_ = x
	<-c
}
