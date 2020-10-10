package main

func main() {
	c := make(chan bool)
	var x interface{}
	go func() {
		x = nil
		c <- true
	}()
	_ = x
	<-c
}
