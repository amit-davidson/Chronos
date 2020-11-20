package main

var GlobalX, GlobalY int = 0, 0
var GlobalCh chan int = make(chan int, 2)

func GlobalFunc1() {
	GlobalY = GlobalX
	GlobalCh <- 1
}

func GlobalFunc2() {
	GlobalX = 1
	GlobalCh <- 1
}

func main() {
	go GlobalFunc1()
	go GlobalFunc2()
	<-GlobalCh
	<-GlobalCh
}
