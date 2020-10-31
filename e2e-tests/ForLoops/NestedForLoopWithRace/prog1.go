package main

func main() {
	m := make(map[string]string)
	go func() {
		for _, _ = range m {
			for _, _ = range m {
				m["1"] = "a"
			}
		}
	}()
	m["2"] = "b"
}
