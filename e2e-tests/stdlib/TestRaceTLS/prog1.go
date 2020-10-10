package main

func main() {
	comm := make(chan *int)
	done := make(chan bool, 2)
	go func() {
		var x int
		comm <- &x
		x = 1
		x = *(<-comm)
		done <- true
	}()
	go func() {
		p := <-comm
		*p = 2
		comm <- p
		done <- true
	}()
	<-done
	<-done
}
