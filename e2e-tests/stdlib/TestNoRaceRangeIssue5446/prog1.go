package main

func main() {
	ch := make(chan int, 3)
	a := []int{1, 2, 3}
	b := []int{4}
	// used to insert a spurious instrumentation of a[i]
	// and crash.
	i := 1
	for i, a[i] = range b {
		ch <- i
	}
	close(ch)
}
