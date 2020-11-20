package main

func main() {
	var x, y int
	_ = y
	ch := make(chan int, 2)

	go func() {
		y = x
		ch <- 1
	}()
	go func() {
		x = 1
		ch <- 1
	}()
	<-ch
	<-ch
}
