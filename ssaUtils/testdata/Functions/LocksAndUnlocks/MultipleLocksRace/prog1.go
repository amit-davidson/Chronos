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
		x = 1
		mu.Lock()
		defer mu.Unlock()
		ch <- true
	}()
	go func() {
		x = 2
		mu.Lock()
		mu.Unlock()
		ch <- true
	}()
	<-ch
	<-ch
}
