package main

func main() {
	const P = 4
	const N = 1e6
	var tinySink *byte
	_ = tinySink
	done := make(chan bool)
	for p := 0; p < P; p++ {
		go func() {
			for i := 0; i < N; i++ {
				var b byte
				if b != 0 {
					tinySink = &b // make it heap allocated
				}
				b = 42
			}
			done <- true
		}()
	}
	for p := 0; p < P; p++ {
		<-done
	}
}
