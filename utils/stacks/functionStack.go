package stacks

import "golang.org/x/tools/go/ssa"

type FunctionStackWithMap struct {
	stack       FunctionStack
	FunctionMap FunctionMap
}

func NewFunctionStackWithMap() *FunctionStackWithMap {
	stack := make([]*ssa.Function, 0)
	FunctionMap := make(FunctionMap)
	basicBlockStack := &FunctionStackWithMap{stack: stack, FunctionMap: FunctionMap}
	return basicBlockStack
}

func (s *FunctionStackWithMap) Copy() *FunctionStackWithMap {
	tmp := NewFunctionStackWithMap()
	tmpStack := s.stack.Copy()
	tmpMap := make(FunctionMap)
	for k, v := range s.FunctionMap {
		tmpMap[k] = v
	}
	tmp.stack = *tmpStack
	tmp.FunctionMap = tmpMap
	return tmp
}

func (s *FunctionStackWithMap) GetItems() *FunctionStack {
	return &s.stack
}

func (s *FunctionStackWithMap) Iter() []*ssa.Function {
	return s.stack
}

func (s *FunctionStackWithMap) Contains(v *ssa.Function) bool {
	_, ok := s.FunctionMap[v]
	return ok
}

func (s *FunctionStackWithMap) Push(v *ssa.Function) {
	s.stack.Push(v)
	s.FunctionMap[v] = struct{}{}
}

func (s *FunctionStackWithMap) Pop() *ssa.Function {
	v := s.stack.Pop()
	delete(s.FunctionMap, v)
	return v
}

func (s *FunctionStackWithMap) Merge(sn *FunctionStackWithMap) {
	for _, item := range sn.Iter() {
		s.stack.Push(item)
		s.FunctionMap[item] = struct{}{}
	}
}

type FunctionMap map[*ssa.Function]struct{}

type FunctionStack []*ssa.Function

func (s *FunctionStack) GetItems() []*ssa.Function {
	return *s
}

func (s *FunctionStack) Copy() *FunctionStack {
	tmp := make([]*ssa.Function, len(*s))
	copy(tmp, *s)
	return (*FunctionStack)(&tmp)
}

func (s *FunctionStack) Push(v *ssa.Function) {
	*s = append(*s, v)
}

func (s *FunctionStack) Pop() *ssa.Function {
	if len(*s) == 0 {
		return nil
	}
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *FunctionStack) MergeStacks(items *FunctionStack) {
	for _, item := range items.GetItems() {
		s.Push(item)
	}
}
