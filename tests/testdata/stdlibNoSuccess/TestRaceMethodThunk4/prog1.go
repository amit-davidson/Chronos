package main

type Base int

func (b *Base) Foo() int {
	return 42
}

func (b Base) Bar() int {
	return int(b)
}

func main() {
	type Derived struct {
		pad int
		*Base
	}
	var d Derived
	d.Base = new(Base)
	done := make(chan bool)
	go func() {
		_ = d.Bar()
		done <- true
	}()
	*(*int)(d.Base) = 42
	<-done
}
