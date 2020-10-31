package main

func main() {
	var x int
	_ = x
	compl := make(chan bool)
	c := make(chan bool, 10)
	c1 := make(chan bool)
	go func() {
		x = 1
		select {
		case c <- true:
		case <-c1:
		}
		compl <- true
	}()
	<-c
	x = 2
	<-compl
}
