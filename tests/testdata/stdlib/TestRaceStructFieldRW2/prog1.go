package main

type Point struct {
	x, y int
}

func main() {
	p := &Point{0, 0}
	ch := make(chan bool, 1)
	go func() {
		p.x = 1
		ch <- true
	}()
	_ = p.x
	<-ch
}
