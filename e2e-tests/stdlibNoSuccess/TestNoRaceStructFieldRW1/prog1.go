package main

type Point struct {
	x, y int
}

type NamedPoint struct {
	name string
	p    Point
}

func main() {
	// Same struct, different variables, no
	// pointers. The layout is known (at compile time?) ->
	// no read on p
	// writes on x and y
	p := Point{0, 0}
	ch := make(chan bool, 1)
	go func() {
		p.x = 1
		ch <- true
	}()
	p.y = 1
	<-ch
	_ = p
}
