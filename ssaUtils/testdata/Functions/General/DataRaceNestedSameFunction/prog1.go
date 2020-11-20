package main

var count int

func race() {
	count++
}

func f() {
	race()
}

func main() {
	go f()
	go f()
}
