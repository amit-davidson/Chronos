package main

type Point struct {
	x, y int
}

func main() {
	var y int
	for x := -1; x <= 1; x++ {
		switch {
		case x < 0:
			y = -1
		case x == 0:
			y = 0
		case x > 0:
			y = 1
		}
	}
	y++
}
