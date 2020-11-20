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
	i := &InterImpl{}
	go func() {
		i = &InterImpl{}
		c <- true
	}()
	i.Foo(0)
	<-c
}
