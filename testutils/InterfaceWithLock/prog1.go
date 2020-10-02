package main

import "sync"

type IceCreamMaker interface {
	// Hello greets a customer
	Hello()
}

type Ben struct {
	name string
}

type Jerry struct {
	name string
}

func (b *Ben) Hello() {
	mutex.Lock()
	b.name = "Ben"
}

func (j *Jerry) Hello() {
	mutex.Lock()
	j.name = "Jerry"
}

var mutex = sync.Mutex{}

func main() {
	var ben = &Ben{}
	var maker IceCreamMaker = ben
	maker.Hello()
	ben.name = "1"
}
