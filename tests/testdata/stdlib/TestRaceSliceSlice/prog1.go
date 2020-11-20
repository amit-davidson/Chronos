package main

func main() {
	c := make(chan bool, 1)
	x := make([]int, 10)
	go func() {
		x = make([]int, 20)
		c <- true
	}()
	_ = x[2:3]
	<-c
}
