package main

func main() {
	ch := make(chan bool, 1)
	var x string
	p := &x
	go func() {
		*p = "a"
		ch <- true
	}()
	_ = *p
	<-ch
}
