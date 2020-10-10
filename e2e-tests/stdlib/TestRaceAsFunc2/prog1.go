package main

func main() {
	c := make(chan bool, 1)
	x := 0
	go func() {
		func(x int) {
		}(x)
		c <- true
	}()
	x = 42
	<-c
}
