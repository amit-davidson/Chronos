package main

import (
	"sync"
)

type SystemAgent struct {
	ServerContext map[string]int `json:"servercontext,omitempty"`
}

var mutex1 = sync.Mutex{}

func main() {
	mutex1.Lock()
	sa := &SystemAgent{}
	sa.ServerContext = map[string]int{"6": 6}
	fn2(sa)
}

//go:noinline
func fn2(item *SystemAgent) {
	item.ServerContext = map[string]int{"7": 7}
	mutex1.Unlock()
}
