package utils

type Stack []int

func NewStack() *Stack {
	stack := make([]int, 0)
	return (*Stack)(&stack)
}

func (s *Stack) GetAllItems() []int {
	tmp := make([]int, len(*s))
	copy(tmp, *s)
	return tmp
}

func (s *Stack) Push(v int) {
	*s = append(*s, v)
}

func (s *Stack) Pop() int {
	if len(*s) == 0 {
		return -1
	}
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}
