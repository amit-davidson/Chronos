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
		switch {
		default:
			x = 1
		case x == 100:
			x = -x
		}
		ch <- 1
	}()
	<-ch
	<-ch
}
