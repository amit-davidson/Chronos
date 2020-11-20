package main

func main() {
	var x int
	_ = x
	compl := make(chan bool)
	c := make(chan bool)
	c1 := make(chan bool)

	go func() {
		x = 1
		// At least two channels are needed because
		// otherwise the compiler optimizes select out.
		// See comment in runtime/select.go:^func selectgo.
		select {
		case c <- true:
		case c1 <- true:
		}
		compl <- true
	}()
	select {
	case <-c:
	case c1 <- true:
	}
	x = 2
	<-compl
}
