package main

import "runtime"

func main() {
	const N = 2
	var a [N]int
	var x, y int
	_ = x + y
	done := make(chan bool, N)
	for i, v := range a {
		go func(i int) {
			// we don't want a write-vs-write race
			// so there is no array b here
			if i == 0 {
				x = v
			} else {
				y = v
			}
			done <- true
		}(i)
		// Ensure the goroutine runs before we continue the loop.
		runtime.Gosched()
	}
	for i := 0; i < N; i++ {
		<-done
	}
}
