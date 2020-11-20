package main

func main() {
	var x, y int
	var i interface{} = &x
	c := make(chan int, 1)
	go func() {
		switch i := i.(type) {
		case nil:
		case *int:
			*i = y
		}
		c <- 1
	}()
	x = y
	<-c
}
