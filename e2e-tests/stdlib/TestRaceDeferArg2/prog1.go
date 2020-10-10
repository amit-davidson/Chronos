package main

type DeferT int

func (d DeferT) Foo() {
}

func main() {
	c := make(chan bool, 1)
	var x DeferT
	go func() {
		var y DeferT
		x = y
		c <- true
	}()
	func() {
		defer x.Foo()
	}()
	<-c
}
