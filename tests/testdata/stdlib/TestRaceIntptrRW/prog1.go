package main

func main() {
	var x, y int
	var p *int = &x
	ch := make(chan bool, 1)
	go func() {
		*p = 5
		ch <- true
	}()
	y = *p
	x = y
	<-ch
}
