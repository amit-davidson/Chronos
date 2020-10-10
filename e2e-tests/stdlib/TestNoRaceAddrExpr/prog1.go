package main

func main() {
	c := make(chan bool, 1)
	x := 0
	go func() {
		x = 42
		c <- true
	}()
	_ = &x
	<-c
}
