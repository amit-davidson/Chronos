package pkg

import (
	"sync"
)

func main() {
	var mu sync.Mutex
	var x int16 = 0
	_ = x
	ch := make(chan bool, 2)
	go func() {
		mu.Lock()
		defer mu.Unlock()
		x = 1
		ch <- true
	}()
	go func() {
		mu.Lock()
		x = 2
		mu.Unlock()
		ch <- true
	}()
	<-ch
	<-ch
}
