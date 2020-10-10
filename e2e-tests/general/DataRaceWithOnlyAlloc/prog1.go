package testutils

import (
	"fmt"
)

func main() {
	fmt.Println(getNumber())
}

func getNumber() int {
	var i int
	go func() {
		i = 5
	}()

	return i
}
