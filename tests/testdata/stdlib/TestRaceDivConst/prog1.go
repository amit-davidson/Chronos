package main

func main() {
	var x, y, z uint32
	_ = x
	ch := make(chan int, 2)

	go func() {
		x = y / 3 // involves only a HMUL node
		ch <- 1
	}()
	go func() {
		y = z
		ch <- 1
	}()
	<-ch
	<-ch
}
