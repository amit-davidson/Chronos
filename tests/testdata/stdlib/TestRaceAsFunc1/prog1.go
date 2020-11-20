package main

func main() {
	var s []byte
	c := make(chan bool, 1)
	go func() {
		var err error
		s, err = func() ([]byte, error) {
			t := []byte("hello world")
			return t, nil
		}()
		c <- true
		_ = err
	}()
	_ = string(s)
	<-c
}
