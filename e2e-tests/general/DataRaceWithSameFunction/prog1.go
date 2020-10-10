package main

import (
	"time"
)

var count int

func race() {
	count++
}

func main() {
	go race()
	go race()
	time.Sleep(1 * time.Second)
}
