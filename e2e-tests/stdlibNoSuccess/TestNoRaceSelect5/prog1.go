package main

import "runtime"

func main() {
	test := func(sel, needSched bool) {
		var x int
		_ = x
		ch := make(chan bool)
		c1 := make(chan bool)

		done := make(chan bool, 2)
		go func() {
			if needSched {
				runtime.Gosched()
			}
			// println(1)
			x = 1
			if sel {
				select {
				case ch <- true:
				case <-c1:
				}
			} else {
				ch <- true
			}
			done <- true
		}()

		go func() {
			// println(2)
			if sel {
				select {
				case <-ch:
				case <-c1:
				}
			} else {
				<-ch
			}
			x = 1
			done <- true
		}()
		<-done
		<-done
	}

	test(true, true)
	test(true, false)
	test(false, true)
	test(false, false)
}
