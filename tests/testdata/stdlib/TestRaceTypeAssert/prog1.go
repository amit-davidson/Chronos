package main

func main() {
	c := make(chan bool, 1)
	x := 0
	var i interface{} = x
	go func() {
		y := 0
		i = y
		c <- true
	}()
	_ = i.(int)
	<-c
}
