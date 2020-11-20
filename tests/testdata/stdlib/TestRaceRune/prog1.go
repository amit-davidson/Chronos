package main

func main() {
	c := make(chan bool)
	var x rune
	go func() {
		x = 1
		c <- true
	}()
	_ = x
	<-c
}
