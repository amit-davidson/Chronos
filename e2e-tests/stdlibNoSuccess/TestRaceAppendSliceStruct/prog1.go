package main

func main() {
	type X struct {
		x, y int
	}
	c := make(chan bool, 1)
	x := make([]X, 10)
	go func() {
		y := make([]X, 0, 10)
		y = append(y, x...)
		c <- true
	}()
	x[1].	y = 42
	<-c
}
