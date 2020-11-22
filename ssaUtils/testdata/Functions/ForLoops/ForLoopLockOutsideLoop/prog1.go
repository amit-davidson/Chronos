package main

import "sync"

func main() {
	m := make(map[string]string)
	m["2"] = "b" // Second conflicting access.
	mutex := sync.Mutex{}
	mutex.Lock()
	for _, _ = range m {
		//for i, j := range m {
		//fmt.Println(i, j)
		//if i > "1" {
		//	fmt.Println(i, j)
		//}
		m["2"] = "c" // Second conflicting access.
		//}
	}
}
