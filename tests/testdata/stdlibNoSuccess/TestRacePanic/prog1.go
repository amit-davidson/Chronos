package main

type Point struct {
	x, y int
}

func main() {
	var x int
	_ = x
	var zero int = 0
	ch := make(chan bool, 2)
	go func() {
		defer func() {
			err := recover()
			if err == nil {
				panic("should be panicking")
			}
			x = 1
			ch <- true
		}()
		var y int = 1 / zero
		zero = y
	}()
	go func() {
		defer func() {
			err := recover()
			if err == nil {
				panic("should be panicking")
			}
			x = 2
			ch <- true
		}()
		var y int = 1 / zero
		zero = y
	}()

	<-ch
	<-ch
	if zero != 0 {
		panic("zero has changed")
	}
}
