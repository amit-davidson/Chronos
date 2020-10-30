package stacks

import "golang.org/x/tools/go/ssa"

type BasicBlockStack struct {
	stack     blocksStack
	blocksMap blocksMap
}

func NewBasicBlockStack() *BasicBlockStack {
	stack := make([]*ssa.BasicBlock, 0)
	blocksMap := make(map[int]*ssa.BasicBlock)
	basicBlockStack := &BasicBlockStack{stack: stack, blocksMap: blocksMap}
	return basicBlockStack
}

func (s *BasicBlockStack) Push(v *ssa.BasicBlock) {
	s.stack.Push(v)
	s.blocksMap[v.Index] = v
}

func (s *BasicBlockStack) Pop() {
	v := s.stack.Pop()
	delete(s.blocksMap, v.Index)
}

func (s *BasicBlockStack) GetAllItems() blocksStack {
	tmp := make(blocksStack, len(s.stack))
	copy(tmp, s.stack)
	return tmp
}

type blocksMap map[int]*ssa.BasicBlock

type blocksStack []*ssa.BasicBlock

func (s *blocksStack) Push(v *ssa.BasicBlock) {
	*s = append(*s, v)
}

func (s *blocksStack) Pop() *ssa.BasicBlock {
	if len(*s) == 0 {
		return nil
	}
	v := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return v
}

func (s *BasicBlockStack) Contains(v *ssa.BasicBlock) bool {
	_, ok := s.blocksMap[v.Index]
	return ok
}
