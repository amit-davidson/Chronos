package stacks

type IntStack []int

func NewIntStack() *IntStack {
	stack := make([]int, 0)
	return (*IntStack)(&stack)
}

func (s *IntStack) GetAllItems() []int {
	tmp := make([]int, len(*s))
	copy(tmp, *s)
	return tmp
}

func (s *IntStack) Push(v int) {
	*s = append(*s, v)
}

func (s *IntStack) Pop() int {
	if len(*s) == 0 {
		return -1
	}
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}
