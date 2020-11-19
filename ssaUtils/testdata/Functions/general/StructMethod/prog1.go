package main

import (
	"fmt"
)

type Ben struct {
	name string
}

func (b *Ben) Hello() {
	fmt.Printf("my name is:%s", b.name)
}

func main() {
	var ben = &Ben{"Ben"}

	go func() {
		ben.name = "Jerry"
	}()
	ben.Hello()
}
