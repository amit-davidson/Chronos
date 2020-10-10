package main

import "sync"

type IceCreamMaker interface {
	// Hello greets a customer
	Hello()
}

type Jerry struct {
	name string
}

func (j *Jerry) Hello() {
	mutex.Lock()
	j.name = "Jerry"
}

func sayHello() {
	var maker IceCreamMaker
	maker = <-channel
	maker.Hello()
}


var mutex = sync.Mutex{}
var channel = make(chan IceCreamMaker)

func main() {
	jerry := &Jerry{}
	channel <- jerry
	go sayHello()
	go sayHello()
	jerry.name = "1"
}
