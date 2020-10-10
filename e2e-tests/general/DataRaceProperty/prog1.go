package testutils

import (
	"fmt"
	"os"
	"time"
)

type Watchdog struct{ last int64 }

func (w *Watchdog) KeepAlive() {
	w.last = time.Now().UnixNano() // First conflicting access.
}

func (w *Watchdog) Start() {
	go func() {
		time.Sleep(time.Second)
		// Second conflicting access.
		if w.last < time.Now().Add(-10*time.Second).UnixNano() {
			fmt.Println("No keepalives for 10 seconds. Dying.")
			os.Exit(1)
		}
	}()
}

func main() {
	wd := Watchdog{}
	go wd.KeepAlive()
	go wd.Start()
}