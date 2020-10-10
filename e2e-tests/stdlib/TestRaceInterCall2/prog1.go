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
	p := InterImpl{}
	var x Inter = p
	z := 0
	go func() {
		z = 42
		c <- true
	}()
	x.Foo(z)
	<-c
}
