package main

func main() {
	done := make(chan bool)
	c := make(chan bool)
	stop := false
	go func() {
		for {
			_, ok := <-c
			if !ok {
				done <- true
				return
			}
			stop = true
		}
	}()
	for !stop {
		c <- true
	}
	close(c)
	<-done
}
