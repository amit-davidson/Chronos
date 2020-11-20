package main

type Point struct {
	x, y int
}

func main() {
	a := make([]int, 10)
	ch := make(chan bool)
	go func() {
		_ = append(a, 1)
		ch <- true
	}()
	a[0] = 1
	<-ch
}
