package main

func main() {
	var x, y int
	var i interface{} = x
	c := make(chan int, 1)
	go func() {
		switch i.(type) {
		case nil:
		case int:
		}
		c <- 1
	}()
	i = y
	<-c
}
