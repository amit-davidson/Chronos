package main

func main() {
	ch := make(chan bool, 1)
	var err error
	go func() {
		err = nil
		ch <- true
	}()
	_ = err
	<-ch
}
