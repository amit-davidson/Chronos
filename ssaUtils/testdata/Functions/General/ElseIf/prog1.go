package Lock

import (
	"math/rand"
)

var a int

func main() {
	a := rand.Int()
	if a > 0 {
		a = 3
	} else if a == 0 {
		a = 4
	} else {
		a = 5
	}
}
