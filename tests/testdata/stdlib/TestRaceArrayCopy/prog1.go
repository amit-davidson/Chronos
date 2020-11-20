package main

func main() {
	ch := make(chan bool, 1)
	var a [5]int
	go func() {
		a[3] = 1
		ch <- true
	}()
	a = [5]int{1, 2, 3, 4, 5}
	<-ch
}
