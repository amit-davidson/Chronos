package main


func main() {
	var x, y, z int
	_ = y
	ch := make(chan int, 2)

	go func() {
		y = x + z
		ch <- 1
	}()
	go func() {
		y = x + z + z
		ch <- 1
	}()
	<-ch
	<-ch
}
