package main

func main() {
	ch := make(chan bool, 1)
	s := ""
	go func() {
		s = "abacaba"
		ch <- true
	}()
	_ = s
	<-ch
}
