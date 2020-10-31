package main

func main() {
	c := make(chan int)
	x := 0
	go func() {
		c <- x
	}()
	for x = 42; false; {
	}
	<-c
}
