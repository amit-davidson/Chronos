package domain

import "github.com/amit-davidson/Chronos/utils/stacks"

type BlockState struct {
	GuardedAccesses   []*GuardedAccess
	DeferredFunctions *stacks.CallCommonStack
}

func GetEmptyBlockState() *BlockState {
	return &BlockState{
		GuardedAccesses:   make([]*GuardedAccess, 0),
		DeferredFunctions: stacks.NewCallCommonStack(),
	}
}

func CreateBlockState(ga []*GuardedAccess, df *stacks.CallCommonStack) *BlockState {
	return &BlockState{
		GuardedAccesses:   ga,
		DeferredFunctions: df,
	}
}

// AddFunctionCallState is used to add the state of a function call to the blocks total state when iterating through it.
// shouldMergeLockset is used depending if the call was using a goroutine or not.
func (existingBlock *BlockState) AddFunctionCallState(newBlock *BlockState) {
	existingBlock.GuardedAccesses = append(existingBlock.GuardedAccesses, newBlock.GuardedAccesses...)
}

// MergeChildBlock merges child block with it's parent in append-like fashion.
// A -> B
// Will Merge B unto A
func (existingBlock *BlockState) MergeChildBlock(newBlock *BlockState) {
	existingBlock.GuardedAccesses = append(existingBlock.GuardedAccesses, newBlock.GuardedAccesses...)
	existingBlock.DeferredFunctions.MergeStacks(newBlock.DeferredFunctions)
}

// MergeSiblingBlock merges sibling blocks in merge-like fashion.
// A -> B
//   -> C
// Will Merge B and C
func (existingBlock *BlockState) MergeSiblingBlock(newBlock *BlockState) {
	existingGAs := make(map[int]*GuardedAccess, len(existingBlock.GuardedAccesses))
	for _, ga := range existingBlock.GuardedAccesses {
		existingGAs[ga.ID] = ga
	}

	for _, newGA := range newBlock.GuardedAccesses {
		if _, ok := existingGAs[newGA.ID]; !ok {
			existingBlock.GuardedAccesses = append(existingBlock.GuardedAccesses, newGA)
		}
	}
}

func (existingBlock *BlockState) Copy() *BlockState {
	newFunctionState := &BlockState{}
	for _, ga := range existingBlock.GuardedAccesses {
		newFunctionState.GuardedAccesses = append(newFunctionState.GuardedAccesses, ga.Copy())
	}
	newFunctionState.DeferredFunctions = existingBlock.DeferredFunctions
	return newFunctionState
}
