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
		case <-c:
		case <-c1:
		}
		compl <- true
	}()
	close(c)
	x = 2
	<-compl
}
