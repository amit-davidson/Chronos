package main

type P struct {
	x, y int
}

type S struct {
	s1, s2 P
}

func main() {
	c := make(chan bool, 1)
	var s S
	go func() {
		s.s1.x = 1
		c <- true
	}()
	s = S{}
	<-c
}
