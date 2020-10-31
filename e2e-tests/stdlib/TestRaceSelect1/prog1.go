package main

func main() {
	var x int
	_ = x
	compl := make(chan bool, 2)
	c := make(chan bool)
	c1 := make(chan bool)

	go func() {
		<-c
		<-c
	}()
	f := func() {
		select {
		case c <- true:
		case c1 <- true:
		}
		x = 1
		compl <- true
	}
	go f()
	go f()
	<-compl
	<-compl
}
