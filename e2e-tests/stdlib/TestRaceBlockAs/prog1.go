package main

func main() {
	c := make(chan bool, 1)
	var x, y int
	go func() {
		x = 42
		c <- true
	}()
	x, y = y, x
	<-c
}
