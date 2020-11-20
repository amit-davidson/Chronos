package main

type DeferT int

func (d DeferT) Foo() {
}

func main() {
	c := make(chan bool, 1)
	x := 0
	go func() {
		x = 42
		c <- true
	}()
	func() {
		defer func(x int) {
		}(x)
	}()
	<-c
}
