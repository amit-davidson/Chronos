package main

type OsFile struct{}

func (*OsFile) Read() {
}

type IoReader interface {
	Read()
}

func main() {
	c := make(chan bool)
	f := &OsFile{}
	go func() {
		go func(x IoReader) {
		}(f)
		c <- true
	}()
	f = &OsFile{}
	<-c
}
