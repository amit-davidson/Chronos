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
	i := 0
	for x = 42; i < 10; i++ {
		c <- true
	}
	close(c)
	<-done
}
