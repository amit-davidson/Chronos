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
	var a, b Writer
	a = DummyWriter{1}
	ch := make(chan bool, 1)
	go func() {
		a = DummyWriter{1}
		ch <- true
	}()
	_ = a == b
	<-ch
}
