package main

func main() {
	type Task struct {
		f    func()
		done chan bool
	}

	queue := make(chan Task)
	dummy := make(chan bool)

	go func() {
		for true {
			select {
			case t := <-queue:
				t.f()
				t.done <- true
			}
		}
	}()

	doit := func(f func()) {
		done := make(chan bool, 1)
		select {
		case queue <- Task{f, done}:
		case <-dummy:
		}
		select {
		case <-done:
		case <-dummy:
		}
	}

	var x int
	doit(func() {
		x = 1
	})
	_ = x
}
