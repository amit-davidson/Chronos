package main

import "sync"

func main() {
	g := &sync.WaitGroup{}
	runs := 5
	recur(runs, g)
}

func recur(iter int, g *sync.WaitGroup) {
	if iter <= 0 {
		return
	}
	go recur(iter-1, g)
}
