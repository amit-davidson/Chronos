package main

func main() {
	c := make(chan bool)
	v := 0
	go func() {
		go func(x interface{}) {
		}(v)
		c <- true
	}()
	v = 42
	<-c
}
