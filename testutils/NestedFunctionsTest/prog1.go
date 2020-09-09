package main

import (
	"sync"
)

type ObjectA struct {
	PropertyA map[string]int `json:"peropertya,omitempty"`
}

var mutex1 = sync.Mutex{}

func main() {
	mutex1.Lock()
	sa := &ObjectA{}
	sa.PropertyA = map[string]int{"6": 6}
	go fn2(sa)
}

//go:noinline
func fn2(item *ObjectA) {
	item.PropertyA = map[string]int{"7": 7}
	mutex1.Unlock()
}
