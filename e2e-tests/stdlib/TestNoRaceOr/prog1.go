package main

func main() {
	c := make(chan bool)
	x, y := 0, 0
	go func() {
		x = 1
		c <- true
	}()
	if y == 0 || x == 1 {
	}
	<-c
}
