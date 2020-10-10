package main

func main() {
	m := make(map[int]map[int]int)
	ch := make(chan int, 1)
	m[0] = make(map[int]int)
	go func() {
		_ = len(m[0])
		ch <- 0
	}()
	m[0][0] = 1
	<-ch
}
