package domain

import "github.com/amit-davidson/Chronos/utils/stacks"

type BlockState struct {
	GuardedAccesses   []*GuardedAccess
	Lockset           *Lockset
	DeferredFunctions *stacks.CallCommonStack
}

func GetEmptyBlockState() *BlockState {
	return &BlockState{
		GuardedAccesses:   make([]*GuardedAccess, 0),
		Lockset:           NewLockset(),
		DeferredFunctions: stacks.NewCallCommonStack(),
	}
}

func CreateBlockState(ga []*GuardedAccess, ls *Lockset, df *stacks.CallCommonStack) *BlockState {
	return &BlockState{
		GuardedAccesses:   ga,
		Lockset:           ls,
		DeferredFunctions: df,
	}
}

// AddFunctionCallState is used to add the state of a function call to the blocks total state when iterating through it.
// shouldMergeLockset is used depending if the call was using a goroutine or not.
func (existingBlock *BlockState) AddFunctionCallState(newBlock *BlockState, shouldMergeLockset bool) {
	for _, guardedAccess := range newBlock.GuardedAccesses {
		guardedAccess.Lockset.UpdateLockSetWithoutCopy(existingBlock.Lockset)

	}
	existingBlock.GuardedAccesses = append(existingBlock.GuardedAccesses, newBlock.GuardedAccesses...)
	if shouldMergeLockset {
		existingBlock.Lockset.UpdateLockSet(newBlock.Lockset.Locks, newBlock.Lockset.Unlocks)
	}
}

// MergeChildBlock merges child block with it's parent in append-like fashion.
// A -> B
// Will Merge B unto A
func (existingBlock *BlockState) MergeChildBlock(newBlock *BlockState) {
	for _, guardedAccess := range newBlock.GuardedAccesses {
		guardedAccess.Lockset.UpdateLockSetWithoutCopy(existingBlock.Lockset)
	}
	existingBlock.GuardedAccesses = append(existingBlock.GuardedAccesses, newBlock.GuardedAccesses...)
	existingBlock.DeferredFunctions.MergeStacks(newBlock.DeferredFunctions)
	existingBlock.Lockset.UpdateLockSet(newBlock.Lockset.Locks, newBlock.Lockset.Unlocks)
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
		if existingGA, ok := existingGAs[newGA.ID]; !ok {
			existingBlock.GuardedAccesses = append(existingBlock.GuardedAccesses, newGA)
		} else {
			existingGA.Lockset.MergeSiblingLockset(newGA.Lockset)
		}
	}

	existingBlock.Lockset.MergeSiblingLockset(newBlock.Lockset)
}

func (existingBlock *BlockState) Copy() *BlockState {
	newFunctionState := GetEmptyBlockState()
	newFunctionState.Lockset = existingBlock.Lockset.Copy()
	for _, ga := range existingBlock.GuardedAccesses {
		newFunctionState.GuardedAccesses = append(newFunctionState.GuardedAccesses, ga.Copy())
	}
	newFunctionState.DeferredFunctions = existingBlock.DeferredFunctions
	return newFunctionState
}
