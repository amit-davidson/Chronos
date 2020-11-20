package main

func main() {
	type Empty struct{}
	type X struct {
		y int64
		Empty
	}
	type Y struct {
		x X
		y int64
	}
	c := make(chan X)
	var y Y
	go func() {
		x := y.x
		c <- x
	}()
	y.y = 42
	<-c
}
