package main

type DummyWriter struct {
	state int
}
type Writer interface {
	Write(p []byte) (n int)
}

func (d DummyWriter) Write(p []byte) (n int) {
	return 0
}

func main() {
	var a, b interface{}
	ch := make(chan bool, 1)
	go func() {
		a = 1
		ch <- true
	}()
	a = 2
	<-ch
	_, _ = a, b
}
