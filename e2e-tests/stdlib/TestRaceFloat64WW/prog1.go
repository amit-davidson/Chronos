package main

func main() {
	var x, y float64
	ch := make(chan bool, 1)
	go func() {
		x = 1.0
		ch <- true
	}()
	x = 2.0
	<-ch

	y = x
	x = y
}
