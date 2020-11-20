package main

func main() {
	var x, y complex128
	ch := make(chan bool, 1)
	go func() {
		x = 2 + 2i
		ch <- true
	}()
	x = 4 + 4i
	<-ch

	y = x
	x = y
}
