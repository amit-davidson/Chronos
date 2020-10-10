package main

func main() {
	c := make(chan bool, 1)
	y := 0
	go func() {
		y = 42
		c <- true
	}()
	x := []int{0, y, 42}
	_ = x
	<-c
}
