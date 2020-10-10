package main

func main() {
	var x, y, z int
	_ = y
	ch := make(chan int, 2)
	z = 1

	go func() {
		y = x
		ch <- 1
	}()
	go func() {
		switch {
		case z == 1:
			fallthrough
		case z == 2:
			x = 2
		}
		ch <- 1
	}()

	<-ch
	<-ch
}
