package stacks

type IntStackWithMap struct {
	stack  IntStack
	intMap IntMap
}

func NewIntStackWithMap() *IntStackWithMap {
	stack := make([]int, 0, 20)
	intMap := make(IntMap)
	basicBlockStack := &IntStackWithMap{stack: stack, intMap: intMap}
	return basicBlockStack
}

func NewIntStackWithMapWithParams(stack IntStack, intMap IntMap) *IntStackWithMap {
	basicBlockStack := &IntStackWithMap{stack: stack, intMap: intMap}
	return basicBlockStack
}

func (s *IntStackWithMap) Copy() *IntStackWithMap {
	tmp := &IntStackWithMap{}
	tmpStack := s.stack.Copy()
	tmpMap := make(IntMap, len(s.stack))
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

func (s *IntStackWithMap) Iter() []int {
	return s.stack
}

func (s *IntStackWithMap) Contains(v int) bool {
	_, ok := s.intMap[v]
	return ok
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

func (s *IntStackWithMap) Merge(sn *IntStackWithMap) {
	for _, item := range sn.Iter() {
		s.stack.Push(item)
	}
}

type IntMap map[int]struct{}

type IntStack []int

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
