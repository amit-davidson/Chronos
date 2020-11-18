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

func (existingBlock *BlockState) AddResult(newBlock *BlockState, shouldMergeLockset bool) {
	existingBlock.GuardedAccesses = append(existingBlock.GuardedAccesses, newBlock.GuardedAccesses...)
	if shouldMergeLockset {
		existingBlock.Lockset.UpdateLockSet(newBlock.Lockset.Locks, newBlock.Lockset.Unlocks)
	}
}

// Merge child B unto A:
// A -> B
// Will Merge B unto A
func (existingBlock *BlockState) MergeChildBlock(newBlock *BlockState, shouldUpdateGALockset bool) {
	if shouldUpdateGALockset {
		for _, guardedAccess := range newBlock.GuardedAccesses {
			guardedAccess.Lockset.UpdateLockSet(existingBlock.Lockset.Locks, existingBlock.Lockset.Unlocks)
		}
	}
	existingBlock.GuardedAccesses = append(existingBlock.GuardedAccesses, newBlock.GuardedAccesses...)
	existingBlock.DeferredFunctions.MergeStacks(newBlock.DeferredFunctions)
	existingBlock.Lockset.UpdateLockSet(newBlock.Lockset.Locks, newBlock.Lockset.Unlocks)
}

// Merge state of nodes from branches:
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
