package stacks

import "golang.org/x/tools/go/ssa"

type FunctionStack []*ssa.CallCommon

func NewFunctionStack() *FunctionStack {
	stack := make([]*ssa.CallCommon, 0)
	return (*FunctionStack)(&stack)
}

func (s *FunctionStack) GetAllItems() []*ssa.CallCommon {
	tmp := make([]*ssa.CallCommon, len(*s))
	copy(tmp, *s)
	return tmp
}

func (s *FunctionStack) Push(v *ssa.CallCommon) {
	*s = append(*s, v)
}

func (s *FunctionStack) Len() int {
	return len(*s)
}

func (s *FunctionStack) Pop() *ssa.CallCommon {
	if len(*s) == 0 {
		return nil
	}
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *FunctionStack) MergeStacks(ns *FunctionStack) {
	for _, item := range ns.GetAllItems() {
		s.Push(item)
	}
}
