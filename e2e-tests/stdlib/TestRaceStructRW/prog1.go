package main

type Point struct {
	x, y int
}

type NamedPoint struct {
	name string
	p    Point
}

func main() {
	p := Point{0, 0}
	ch := make(chan bool, 1)
	go func() {
		p = Point{1, 1}
		ch <- true
	}()
	q := p
	<-ch
	p = q
}
