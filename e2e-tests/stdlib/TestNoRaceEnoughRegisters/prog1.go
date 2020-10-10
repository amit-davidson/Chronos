package main


func main() {
	// from erf.go
	const (
		sa1 = 1
		sa2 = 2
		sa3 = 3
		sa4 = 4
		sa5 = 5
		sa6 = 6
		sa7 = 7
		sa8 = 8
	)
	var s, S float64
	s = 3.1415
	S = 1 + s*(sa1+s*(sa2+s*(sa3+s*(sa4+s*(sa5+s*(sa6+s*(sa7+s*sa8)))))))
	s = S
}
