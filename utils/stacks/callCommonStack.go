package stacks

import "golang.org/x/tools/go/ssa"

type CallCommonStack []*ssa.CallCommon

func NewCallCommonStack() *CallCommonStack {
	stack := make([]*ssa.CallCommon, 0)
	return (*CallCommonStack)(&stack)
}

func (s *CallCommonStack) GetItems() []*ssa.CallCommon {
	tmp := make([]*ssa.CallCommon, len(*s))
	copy(tmp, *s)
	return tmp
}

func (s *CallCommonStack) Push(v *ssa.CallCommon) {
	*s = append(*s, v)
}

func (s *CallCommonStack) Len() int {
	return len(*s)
}

func (s *CallCommonStack) Pop() *ssa.CallCommon {
	if len(*s) == 0 {
		return nil
	}
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *CallCommonStack) MergeStacks(ns *CallCommonStack) {
	for _, item := range ns.GetItems() {
		s.Push(item)
	}
}
