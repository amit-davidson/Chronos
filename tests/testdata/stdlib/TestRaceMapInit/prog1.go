package main

func main() {
	c := make(chan bool, 1)
	y := 0
	go func() {
		y = 42
		c <- true
	}()
	x := map[int]int{0: 42, y: 42}
	_ = x
	<-c
}
