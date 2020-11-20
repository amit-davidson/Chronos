package main

func main() {
	c := make(chan bool, 1)
	x := make([]int, 10)
	i := 2
	go func() {
		i = 3
		c <- true
	}()
	_ = x[i:4]
	<-c
}
