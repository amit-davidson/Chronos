package domain

import "github.com/amit-davidson/Chronos/utils/stacks"

type FunctionState struct {
	GuardedAccesses []*GuardedAccess
	Lockset         *Lockset
}

func GetFunctionState() *FunctionState {
	return &FunctionState{
		GuardedAccesses: make([]*GuardedAccess, 0),
		Lockset:         NewLockset(),
	}
}

func CreateFunctionState(ga []*GuardedAccess, ls *Lockset) *FunctionState {
	return &FunctionState{
		GuardedAccesses: ga,
		Lockset:         ls,
	}
}

// Merge child B unto A:
// A -> B
// Will Merge B unto A
func (fs *FunctionState) MergeChildFunction(newFunction *FunctionState, shouldMergeLockset bool) {
	fs.GuardedAccesses = append(fs.GuardedAccesses, newFunction.GuardedAccesses...)
	if shouldMergeLockset {
		fs.Lockset.UpdateLockSet(newFunction.Lockset.ExistingLocks, newFunction.Lockset.ExistingUnlocks)
	}
}

func (fs *FunctionState) UpdateFunctionWithContext(context *Context) {
	for _, ga := range fs.GuardedAccesses {
		ga.ID = context.GuardedAccessCounter.GetNext()
		ga.State.GoroutineID = context.GoroutineID
		context.Increment()
		ga.State.Clock = context.Copy().Clock

		newStack := context.StackTrace.Copy().GetItems()
		gaStack := ga.Stacktrace.GetItems()
		diffPoint := 0
		for i := 0; i < len(gaStack); i++ {
			pos := newStack[i]
			gaPos := gaStack[i]
			if pos != gaPos { // We reached the point where the paths differ
				diffPoint = i
			}
		}

		for i := diffPoint + 1; i < len(gaStack); i++ {
			newStack = append(newStack, gaStack[i])
		}

		ga.Stacktrace = (*stacks.IntStack)(&newStack)
	}
}

func (fs *FunctionState) Copy() *FunctionState {
	newFunctionState := GetFunctionState()
	newFunctionState.Lockset = fs.Lockset.Copy()
	for _, ga := range fs.GuardedAccesses {
		newFunctionState.GuardedAccesses = append(newFunctionState.GuardedAccesses, ga.Copy())
	}
	return newFunctionState
}
