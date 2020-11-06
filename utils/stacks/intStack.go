package stacks

type IntStackWithMap struct {
	stack  IntStack
	intMap IntMap
}

func NewIntStackWithMap() *IntStackWithMap {
	stack := make([]int, 0)
	intMap := make(IntMap)
	basicBlockStack := &IntStackWithMap{stack: stack, intMap: intMap}
	return basicBlockStack
}

func (s *IntStackWithMap) Copy() *IntStackWithMap {
	tmp := NewIntStackWithMap()
	tmpStack := s.stack.Copy()
	tmpMap := make(IntMap)
	for k, v := range s.intMap {
		tmpMap[k] = v
	}
	tmp.stack = *tmpStack
	tmp.intMap = tmpMap
	return tmp
}

func (s *IntStackWithMap) GetItems() *IntStack {
	return &s.stack
}

func (s *IntStackWithMap) Push(v int) {
	s.stack.Push(v)
	s.intMap[v] = struct{}{}
}

func (s *IntStackWithMap) Pop() int {
	v := s.stack.Pop()
	delete(s.intMap, v)
	return v
}

type IntMap map[int]struct{}

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

func (s *IntStack) MergeStacks(items *IntStack) {
	for _, item := range items.GetItems() {
		s.Push(item)
	}
}
