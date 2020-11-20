package main

func main() {
	done := make(chan bool, 1)
	c1 := make(chan bool, 1)
	c2 := make(chan bool)
	var x, y int
	go func() {
		select {
		case c1 <- true:
			x = 1
		case c2 <- true:
			y = 1
		}
		done <- true
	}()
	_ = x
	_ = y
	<-done
}
