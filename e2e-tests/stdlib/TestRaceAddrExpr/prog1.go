package main

type AddrT struct {
	_ [256]byte
	x int
}

type AddrT2 struct {
	_ [512]byte
	p *AddrT
}

func main() {
	c := make(chan bool, 1)
	a := AddrT2{p: &AddrT{x: 42}}
	go func() {
		a.p = &AddrT{x: 43}
		c <- true
	}()
	_ = &a.p.x
	<-c
}
