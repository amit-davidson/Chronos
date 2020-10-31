package main

func main() {
	var x int
	_ = x
	compl := make(chan bool)
	c := make(chan bool)
	c1 := make(chan bool)
	go func() {
		x = 1
		select {
		case c <- true:
		case c1 <- true:
		}
		compl <- true
	}()
	x = 2
	select {
	case <-c:
	}
	<-compl
}
