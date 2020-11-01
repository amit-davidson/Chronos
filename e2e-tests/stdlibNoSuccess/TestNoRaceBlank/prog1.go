package main

func main() {
	var a [5]int
	ch := make(chan bool, 1)
	go func() {
		_, _ = a[0], a[1]
		ch <- true
	}()
	_, _ = a[2], a[3]
	<-ch
	a[1] = a[0]
}
