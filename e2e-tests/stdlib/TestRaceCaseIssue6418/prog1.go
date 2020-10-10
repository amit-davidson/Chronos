package main

func main() {
	m := map[string]map[string]string{
		"a": {
			"b": "c",
		},
	}
	ch := make(chan int)
	go func() {
		m["a"]["x"] = "y"
		ch <- 1
	}()
	switch m["a"]["b"] {
	}
	<-ch
}
