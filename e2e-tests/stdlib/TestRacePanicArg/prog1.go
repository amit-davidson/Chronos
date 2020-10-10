package main

import "errors"

type DeferT int

func (d DeferT) Foo() {
}

func main() {
	c := make(chan bool, 1)
	err := errors.New("err")
	go func() {
		err = errors.New("err2")
		c <- true
	}()
	defer func() {
		recover()
		<-c
	}()
	panic(err)
}
