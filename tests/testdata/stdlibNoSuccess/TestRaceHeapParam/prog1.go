package main

func main() {
	done := make(chan bool)
	x := func() (x int) {
		go func() {
			x = 42
			done <- true
		}()
		return
	}()
	_ = x
	<-done
}
