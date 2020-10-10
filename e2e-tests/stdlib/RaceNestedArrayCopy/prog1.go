package main

type Point struct {
	x, y int
}

func main() {
	ch := make(chan bool, 1)
	type (
		Point32   [2][2][2][2][2]Point
		Point1024 [2][2][2][2][2]Point32
		Point32k  [2][2][2][2][2]Point1024
		Point1M   [2][2][2][2][2]Point32k
	)
	var a, b Point1M
	go func() {
		a[0][1][0][1][0][1][0][1][0][1][0][1][0][1][0][1][0][1][0][1].y = 1
		ch <- true
	}()
	a = b
	<-ch
}
