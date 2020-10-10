package main

func main() {
	c := make(chan bool, 1)
	x := "hello"
	go func() {
		x = "world"
		c <- true
	}()
	_ = x[2:3]
	<-c
}
