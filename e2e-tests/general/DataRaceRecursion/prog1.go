package main

var a int

func main() {
	go func() {
		a = 5
	}()
	fn1(0)
}

func fn1(counter int) {
	if counter >= 10 {
		return
	}
	if counter == 8 {
		a = 6
	}
	counter += 1
	fn1(counter)
}
