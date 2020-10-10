package main

func main() {
	var x int
	ch := make(chan bool, 2)
	go func() {
		x = 42
		ch <- true
	}()
	go func(y int) {
		ch <- true
	}(x)
	<-ch
	<-ch
}
