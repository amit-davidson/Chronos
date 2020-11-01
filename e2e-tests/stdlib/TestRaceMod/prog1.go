package main

func main() {
	var x, y, z int
	_ = x
	ch := make(chan int, 2)

	go func() {
		x = y % (z + 1)
		ch <- 1
	}()
	go func() {
		y = z
		ch <- 1
	}()
	<-ch
	<-ch
}
