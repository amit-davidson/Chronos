package stacks

import "golang.org/x/tools/go/ssa"

// Used by CFG to traverse the graph. It uses both as a stack for traversal by order, and as a map to count occurrences
// and fast retrieval of items.

type BlockMap map[int]struct{}

func NewBlockMap() *BlockMap {
	blocksMap := make(BlockMap)
	return &blocksMap
}

func (s *BlockMap) Add(v *ssa.BasicBlock) {
	(*s)[v.Index] = struct{}{}
}

func (s *BlockMap) Remove(v *ssa.BasicBlock) {
	delete(*s, v.Index)
}

func (s *BlockMap) Contains(v int) bool {
	_, ok := (*s)[v]
	return ok
}
