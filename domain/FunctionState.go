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

// AddContextToFunction adds flow specific context data
func (fs *FunctionState) AddContextToFunction(context *Context) {
	for _, ga := range fs.GuardedAccesses {
		ga.ID = GuardedAccessCounter.GetNext()
		ga.State.GoroutineID = context.GoroutineID
		context.Increment()

		relativePos := ga.State.StackTrace.Iter()[ga.PosToRemove+1:]
		tmpContext := context.CopyWithoutMap()
		tmpContext.StackTrace.GetItems().MergeStacks((*stacks.IntStack)(&relativePos))
		ga.State.StackTrace = tmpContext.StackTrace
		ga.State.Clock = tmpContext.Clock
	}
}

// RemoveContextFromFunction strips any context related data from the guarded access fields. It nullifies id, goroutine id,
// clock and removes from the guarded access the prefix that matches the context path. This way, other flows can take
// the guarded access and add relevant data.
func (fs *FunctionState) RemoveContextFromFunction() {
	gas := make([]*GuardedAccess, 0, len(fs.GuardedAccesses))
	for i := range fs.GuardedAccesses {
		ga := fs.GuardedAccesses[i].ShallowCopy()
		ga.PosToRemove++
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
