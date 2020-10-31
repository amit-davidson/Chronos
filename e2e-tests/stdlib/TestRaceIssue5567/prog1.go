package main

import (
	"crypto/sha1"
	"io"
	"os"
	"runtime"
)

func main() {
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(4))
	in := make(chan []byte)
	res := make(chan error)
	go func() {
		var err error
		defer func() {
			close(in)
			res <- err
		}()
		path := "mop_test.go"
		f, err := os.Open(path)
		if err != nil {
			return
		}
		defer f.Close()
		var n, total int
		b := make([]byte, 17) // the race is on b buffer
		for err == nil {
			n, err = f.Read(b)
			total += n
			if n > 0 {
				in <- b[:n]
			}
		}
		if err == io.EOF {
			err = nil
		}
	}()
	h := sha1.New()
	for b := range in {
		h.Write(b)
	}
	_ = h.Sum(nil)
	err := <-res
	if err != nil {
		os.Exit(1)
	}
}
