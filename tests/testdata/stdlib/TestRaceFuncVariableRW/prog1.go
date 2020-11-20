package main

func main() {
	var f func(x int) int
	f = func(x int) int {
		return x * x
	}
	ch := make(chan bool, 1)
	go func() {
		f = func(x int) int {
			return x
		}
		ch <- true
	}()
	y := f(1)
	<-ch
	x := y
	y = x
	<-ch
}
