package stacks

type IntStack []int

func NewIntStack() *IntStack {
	stack := make([]int, 0)
	return (*IntStack)(&stack)
}

func (s *IntStack) GetItems() []int {
	return *s
}

func (s *IntStack) Copy() *IntStack {
	tmp := make([]int, len(*s))
	copy(tmp, *s)
	return (*IntStack)(&tmp)
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
