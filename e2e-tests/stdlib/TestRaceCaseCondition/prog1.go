package main

func main() {
	var x int = 0
	ch := make(chan int, 2)

	go func() {
		x = 2
		ch <- 1
	}()
	go func() {
		switch x < 2 {
		case true:
			x = 1
			//case false:
			//	x = 5
		}
		ch <- 1
	}()
	<-ch
	<-ch
}
