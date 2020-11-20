package main

type Inter interface {
	Foo(x int)
}
type InterImpl struct {
	x, y int
}

//go:noinline
func (p InterImpl) Foo(x int) {
}

type InterImpl2 InterImpl

func (p *InterImpl2) Foo(x int) {
	if p == nil {
		InterImpl{}.Foo(x)
	}
	InterImpl(*p).Foo(x)
}

func main() {
	c := make(chan bool, 1)
	f := func(x, y int) {}
	x, y := 0, 0
	go func() {
		y = 42
		c <- true
	}()
	f(x, y)
	<-c
}
