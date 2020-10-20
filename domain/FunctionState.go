package domain

type FunctionState struct {
	GuardedAccesses   []*GuardedAccess
	Lockset           *Lockset
	DeferredFunctions []*DeferFunction
}

func GetEmptyFunctionState() *FunctionState {
	return &FunctionState{
		GuardedAccesses:   make([]*GuardedAccess, 0),
		Lockset:           NewLockset(),
		DeferredFunctions: make([]*DeferFunction, 0),
	}
}

func (funcState *FunctionState) MergeStates(funcStateToMerge *FunctionState, shouldMergeLockset bool) {
	funcState.GuardedAccesses = append(funcState.GuardedAccesses, funcStateToMerge.GuardedAccesses...)
	funcState.DeferredFunctions = append(funcState.DeferredFunctions, funcStateToMerge.DeferredFunctions...)
	if shouldMergeLockset {
		funcState.Lockset.UpdateLockSet(funcStateToMerge.Lockset.ExistingLocks, funcStateToMerge.Lockset.ExistingUnlocks)
	}
}

func (funcState *FunctionState) UpdateGuardedAccessesWithLockset(prevLockset *Lockset) {
	for _, guardedAccess := range funcState.GuardedAccesses {
		tempLockset := prevLockset.Copy()
		tempLockset.UpdateLockSet(guardedAccess.Lockset.ExistingLocks, guardedAccess.Lockset.ExistingUnlocks)
		guardedAccess.Lockset = tempLockset
	}
}

func (funcState *FunctionState) MergeBranchState(funcStateToMerge *FunctionState) {
	funcState.RemoveDuplicateGuardedAccess(funcStateToMerge.GuardedAccesses)
	funcState.DeferredFunctions = append(funcState.DeferredFunctions, funcStateToMerge.DeferredFunctions...)
	funcState.Lockset.MergeBranchesLockset(funcStateToMerge.Lockset)
}

func (funcState *FunctionState) RemoveDuplicateGuardedAccess(GuardedAccesses []*GuardedAccess) {
	for _, guardedAccessA := range GuardedAccesses {
		if !contains(funcState.GuardedAccesses, guardedAccessA) {
			funcState.GuardedAccesses = append(funcState.GuardedAccesses, guardedAccessA)
		}
	}
}

func (funcState *FunctionState) Copy() *FunctionState {
	newFunctionState := GetEmptyFunctionState()
	newFunctionState.Lockset = funcState.Lockset.Copy()
	newFunctionState.GuardedAccesses = append(newFunctionState.GuardedAccesses, funcState.GuardedAccesses...)
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
