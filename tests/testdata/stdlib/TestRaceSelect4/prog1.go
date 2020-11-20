package main

func main() {
	done := make(chan bool, 1)
	var x int
	go func() {
		select {
		default:
			x = 2
		}
		done <- true
	}()
	_ = x
	<-done
}
