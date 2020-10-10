package main

func main() {
	var x, y int
	m := make(map[int]map[int]interface{})
	m[0] = make(map[int]interface{})
	c := make(chan int, 1)
	go func() {
		switch i := m[0][1].(type) {
		case nil:
		case *int:
			*i = x
		}
		c <- 1
	}()
	m[0][1] = y
	<-c
}
