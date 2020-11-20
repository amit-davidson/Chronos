package main

func main() {
	done := make(chan bool)
	c := make(chan bool)
	x := 0
	go func() {
		for {
			_, ok := <-c
			if !ok {
				done <- true
				return
			}
			x++
		}
	}()
	for i := 0; i < 10; x++ {
		i++
		c <- true
	}
	close(c)
	<-done
}
