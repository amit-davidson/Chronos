package main

func main() {
	type X struct {
		x, y int
	}
	type Y struct {
		x X
	}
	c := make(chan Y)
	var y Y
	go func() {
		c <- y
	}()
	y.x.y = 42
	<-c
}
