package main

type Point struct {
	x, y int
}

func main() {
	a := make([]int, 0)
	ch := make(chan string)
	go func() {
		a = append(a, 1)
		ch <- ""
	}()
	_ = cap(a)
	<-ch
}
