package main

func main() {
	var f func(x int) int
	_ = f
	ch := make(chan bool, 1)
	go func() {
		f = func(x int) int {
			return x
		}
		ch <- true
	}()
	f = func(x int) int {
		return x * x
	}
	<-ch
}
