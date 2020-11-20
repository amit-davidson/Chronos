package main

type A struct {
	X int
}

func main() {
	a := A{}
	go func() {
		a.X = 1
	}()
	a.X = 2
}
