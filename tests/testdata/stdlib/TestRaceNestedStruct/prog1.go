package main

func main() {
	type X struct {
		x, y int
	}
	type Y struct {
		x X
	}
	c := make(chan Y)
	var t Y
	go func() {
		c <- t
	}()
	//time.Sleep(2*time.Second)
	t.x.y = 42
	val := <-c
	_ = val
}
