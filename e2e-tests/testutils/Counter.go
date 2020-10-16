package testutils

var counter = 0

func GetCounter() int {
	counter += 1
	return counter
}
