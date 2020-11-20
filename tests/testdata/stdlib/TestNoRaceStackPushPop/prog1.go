package main

type stack []int

func (s *stack) push(x int) {
	*s = append(*s, x)
}

func (s *stack) pop() int {
	i := len(*s)
	n := (*s)[i-1]
	*s = (*s)[:i-1]
	return n
}

func main() {
	var s stack
	go func(s *stack) {}(&s)
	s.push(1)
	x := s.pop()
	_ = x
}
