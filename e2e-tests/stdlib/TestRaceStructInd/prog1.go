package main

func main() {
	c := make(chan bool, 1)
	type Item struct {
		x, y int
	}
	i := Item{}
	go func(p *Item) {
		*p = Item{}
		c <- true
	}(&i)
	i.y = 42
	<-c
}
