package domain

import "golang.org/x/tools/go/ssa"

type FunctionState struct {
	GuardedAccesses   []*GuardedAccess
	Lockset           *Lockset
	DeferredFunctions []*ssa.CallCommon
}

func GetEmptyFunctionState() *FunctionState {
	return &FunctionState{
		GuardedAccesses:   make([]*GuardedAccess, 0),
		Lockset:           NewEmptyLockSet(),
		DeferredFunctions: make([]*ssa.CallCommon, 0),
	}
}

func (funcState *FunctionState) MergeStates(funcStateToMerge *FunctionState) {
	funcState.GuardedAccesses = append(funcState.GuardedAccesses, funcStateToMerge.GuardedAccesses...)
	funcState.DeferredFunctions = append(funcState.DeferredFunctions, funcStateToMerge.DeferredFunctions...)
	funcState.Lockset.UpdateLockSet(funcStateToMerge.Lockset.ExistingLocks, funcStateToMerge.Lockset.ExistingUnlocks)
}

func (funcState *FunctionState) MergeStatesAfterGoroutine(funcStateToMerge *FunctionState) {
	funcState.GuardedAccesses = append(funcState.GuardedAccesses, funcStateToMerge.GuardedAccesses...)
	funcState.DeferredFunctions = append(funcState.DeferredFunctions, funcStateToMerge.DeferredFunctions...)
	//funcState.Lockset.UpdateLockSet(funcStateToMerge.Lockset.ExistingLocks, funcStateToMerge.Lockset.ExistingUnlocks)
}

func (funcState *FunctionState) UpdateGuardedAccessesWithLockset(prevLockset *Lockset) {
	for _, guardedAccess := range funcState.GuardedAccesses {
		tempLockset := prevLockset.Copy()
		tempLockset.UpdateLockSet(guardedAccess.Lockset.ExistingLocks, guardedAccess.Lockset.ExistingUnlocks)
		guardedAccess.Lockset = tempLockset
	}
}

func (funcState *FunctionState) MergeBlockStates(funcStateToMerge *FunctionState) {
	funcState.MergeBlockGuardedAccess(funcStateToMerge.GuardedAccesses)
	funcState.DeferredFunctions = append(funcState.DeferredFunctions, funcStateToMerge.DeferredFunctions...)
	funcState.Lockset.MergeBlockLockset(funcStateToMerge.Lockset)
}

func (funcState *FunctionState) MergeBlockGuardedAccess(GuardedAccesses []*GuardedAccess) {
	for _, guardedAccessA := range GuardedAccesses {
		if !contains(funcState.GuardedAccesses, guardedAccessA) {
			funcState.GuardedAccesses = append(funcState.GuardedAccesses, guardedAccessA)
		}
	}
}

func (funcState *FunctionState) Copy() *FunctionState {
	newFunctionState := GetEmptyFunctionState()
	newFunctionState.Lockset = funcState.Lockset.Copy()
	for _, guardedAccessToCopy := range funcState.GuardedAccesses {
		newFunctionState.GuardedAccesses = append(newFunctionState.GuardedAccesses, guardedAccessToCopy)
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
