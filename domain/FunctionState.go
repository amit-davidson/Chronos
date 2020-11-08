package domain

import (
	"github.com/amit-davidson/Chronos/utils"
	"github.com/amit-davidson/Chronos/utils/stacks"
)

var GoroutineCounter *utils.Counter
var GuardedAccessCounter *utils.Counter
var PosIDCounter *utils.Counter


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
		fs.Lockset.UpdateLockSet(newFunction.Lockset.Locks, newFunction.Lockset.Unlocks)
	}
}

// AddContextToFunction adds flow specific context data
func (fs *FunctionState) AddContextToFunction(context *Context) {
	for _, ga := range fs.GuardedAccesses {
		ga.ID = GuardedAccessCounter.GetNext()
		ga.State.GoroutineID = context.GoroutineID
		context.Increment()
		ga.State.Clock = context.Copy().Clock

		tmpStack := context.StackTrace.Copy()
		tmpStack.Merge(ga.State.StackTrace)
		ga.State.StackTrace = tmpStack
	}
}

// RemoveContextFromFunction strips any context related data from the guarded access fields. It nullifies id, goroutine id,
// clock and removes from the guarded access the prefix that matches the context path. This way, other flows can take
// the guarded access and add relevant data.
func (fs *FunctionState) RemoveContextFromFunction(context *Context) {
	gas := make([]*GuardedAccess, 0, len(fs.GuardedAccesses))
	for i := range fs.GuardedAccesses {
		ga := fs.GuardedAccesses[i].Copy()
		ga.ID = 0
		ga.State.GoroutineID = 0
		ga.State.Clock = nil

		newStack := context.StackTrace.Iter()
		gaStack := ga.State.StackTrace.Iter()
		diffPoint := len(newStack)
		for i := 0; i < len(newStack); i++ {
			pos := newStack[i]
			gaPos := gaStack[i]
			if pos != gaPos { // We reached the point where the paths differ
				diffPoint = i
			}
		}

		tmpStack := stacks.NewIntStackWithMap()
		for _, item := range newStack[diffPoint:] {
			tmpStack.Push(item)
		}
		ga.State.StackTrace = tmpStack
		gas = append(gas, ga)
	}
	fs.GuardedAccesses = gas
}

func (fs *FunctionState) Copy() *FunctionState {
	newFunctionState := GetFunctionState()
	newFunctionState.Lockset = fs.Lockset.Copy()
	for _, ga := range fs.GuardedAccesses {
		newFunctionState.GuardedAccesses = append(newFunctionState.GuardedAccesses, ga.Copy())
	}
	return newFunctionState
}
