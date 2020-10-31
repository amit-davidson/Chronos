package domain

import "github.com/amit-davidson/Chronos/utils/stacks"

type FunctionState struct {
	GuardedAccesses   []*GuardedAccess
	Lockset           *Lockset
	DeferredFunctions *stacks.FunctionStack
}

func GetEmptyFunctionState() *FunctionState {
	return &FunctionState{
		GuardedAccesses:   make([]*GuardedAccess, 0),
		Lockset:           NewLockset(),
		DeferredFunctions: stacks.NewFunctionStack(),
	}
}

// Merge state of nodes from the same branch:
// A -> B -> C
// Will Merge B and C
func (funcState *FunctionState) MergeStates(funcStateToMerge *FunctionState, shouldMergeLockset bool) {
	funcState.GuardedAccesses = append(funcState.GuardedAccesses, funcStateToMerge.GuardedAccesses...)
	funcState.DeferredFunctions.MergeStacks(funcStateToMerge.DeferredFunctions)
	if shouldMergeLockset {
		funcState.Lockset.UpdateLockSet(funcStateToMerge.Lockset.ExistingLocks, funcStateToMerge.Lockset.ExistingUnlocks)
	}
}

// Merge state of nodes from branches:
// A -> B
//   -> C
// Will Merge B and C
func (funcState *FunctionState) MergeBranchState(funcStateToMerge *FunctionState) {
	funcState.AppendGuardedAccessWithoutDuplicates(funcStateToMerge.GuardedAccesses)
	for _, mergeGuardedAccess := range funcStateToMerge.GuardedAccesses {
		for _, existingGuardedAccess := range funcState.GuardedAccesses {
			if existingGuardedAccess.ID == mergeGuardedAccess.ID {
				existingGuardedAccess.Lockset.MergeBranchesLockset(mergeGuardedAccess.Lockset)
			}
		}
	}
	funcState.Lockset.MergeBranchesLockset(funcStateToMerge.Lockset)
}

func (funcState *FunctionState) AppendGuardedAccessWithoutDuplicates(GuardedAccesses []*GuardedAccess) {
	for _, guardedAccessA := range GuardedAccesses {
		if !contains(funcState.GuardedAccesses, guardedAccessA) {
			funcState.GuardedAccesses = append(funcState.GuardedAccesses, guardedAccessA)
		}
	}
}

func (funcState *FunctionState) UpdateGuardedAccessesWithLockset(ls *Lockset) {
	for _, guardedAccess := range funcState.GuardedAccesses {
		guardedAccess.Lockset.UpdateLockSet(ls.ExistingLocks, ls.ExistingUnlocks)
	}
}

func (funcState *FunctionState) Copy() *FunctionState {
	newFunctionState := GetEmptyFunctionState()
	newFunctionState.Lockset = funcState.Lockset.Copy()
	for _, ga := range funcState.GuardedAccesses {
		newFunctionState.GuardedAccesses = append(newFunctionState.GuardedAccesses, ga.Copy())
	}
	newFunctionState.DeferredFunctions = funcState.DeferredFunctions
	return newFunctionState
}

func contains(GuardedAccesses []*GuardedAccess, GuardedAccessToCheck *GuardedAccess) bool {
	for _, GuardedAccess := range GuardedAccesses {
		if GuardedAccess.ID == GuardedAccessToCheck.ID {
			return true
		}
	}
	return false
}
