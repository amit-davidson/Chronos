package main

type IceCreamMaker interface {
	Hello()
}

type Ben struct {
	name string
}

func (b *Ben) Hello() {
	b.name = "Ben"
}

func main() {
	var ben = &Ben{}
	var maker IceCreamMaker = ben
	go maker.Hello()
	ben.name = "1"
}
