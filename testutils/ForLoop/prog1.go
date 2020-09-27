package main

import (
	"fmt"
	"sync"
)

func main() {
	m := make(map[string]string)
	m["2"] = "b" // Second conflicting access.
	for k, v := range m {
		fmt.Println(k, v)
		mutex := sync.Mutex{}
		mutex.Lock()
		m["2"] = "b" // Second conflicting access.
	}
}
