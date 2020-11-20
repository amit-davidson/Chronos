package main

type Point struct {
	x, y int
}

type NamedPoint struct {
	name string
	p    Point
}

func main() {
	// Same as NoRaceStructFieldRW1
	// but p is a pointer, so there is a read on p
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
