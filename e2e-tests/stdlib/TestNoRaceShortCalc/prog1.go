package main

func main() {
	c := make(chan bool)
	x, y := 0, 0
	go func() {
		y = 1
		c <- true
	}()
	if x == 0 || y == 0 {
	}
	<-c
}
