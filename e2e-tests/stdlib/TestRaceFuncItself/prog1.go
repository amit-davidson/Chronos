package main

func main() {
	c := make(chan bool)
	f := func() {}
	go func() {
		f()
		c <- true
	}()
	f = func() {}
	<-c
}
