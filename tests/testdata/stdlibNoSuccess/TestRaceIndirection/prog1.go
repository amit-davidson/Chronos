package main

func main() {
	ch := make(chan struct{}, 1)
	var y int
	var x *int = &y
	go func() {
		*x = 1
		ch <- struct{}{}
	}()
	*x = 2
	<-ch
	_ = *x
}
