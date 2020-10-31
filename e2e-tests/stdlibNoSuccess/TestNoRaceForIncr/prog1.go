package main

func main() {
	done := make(chan bool)
	x := 0
	go func() {
		x++
		done <- true
	}()
	for i := 0; i < 0; x++ {
	}
	<-done
}
