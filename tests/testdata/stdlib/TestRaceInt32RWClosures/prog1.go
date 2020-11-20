package main

func main() {
	var x, y int32
	_ = y
	ch := make(chan bool, 2)

	go func() {
		y = x
		ch <- true
	}()
	go func() {
		x = 1
		ch <- true
	}()
	<-ch
	<-ch
}
