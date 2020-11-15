package domain

import (
	"github.com/amit-davidson/Chronos/utils/stacks"
)

// Flow context
type Context struct {
	GoroutineID int
	Clock       VectorClock
	StackTrace  *stacks.IntStackWithMap
}

func NewEmptyContext() *Context {
	return &Context{
		Clock:       VectorClock{},
		GoroutineID: GoroutineCounter.GetNext(),
		StackTrace:  stacks.NewEmptyIntStackWithMap(),
	}
}

func NewGoroutineExecutionState(state *Context) *Context {
	state.Increment()
	return &Context{
		Clock:       state.Clock,
		GoroutineID: GoroutineCounter.GetNext(),
		StackTrace:  state.StackTrace,
	}
}

func (gs *Context) Increment() {
	gs.Clock[gs.GoroutineID] += 1
}

func (gs *Context) MayConcurrent(state *Context) bool {
	timestampAidA := gs.Clock.Get(gs.GoroutineID)
	timestampAidB := state.Clock.Get(gs.GoroutineID)
	timestampBidA := gs.Clock.Get(state.GoroutineID)
	timestampBidB := state.Clock.Get(state.GoroutineID)
	isBefore := timestampAidA <= timestampAidB && timestampBidA < timestampBidB
	isAfter := timestampBidB <= timestampBidA && timestampAidB < timestampAidA
	return !(isBefore || isAfter)
}

func (gs *Context) Copy() *Context {
	return &Context{
		GoroutineID: gs.GoroutineID,
		Clock:       gs.Clock.Copy(),
		StackTrace:  gs.StackTrace.Copy(),
	}
}

func (gs *Context) CopyWithoutMap() *Context {
	return &Context{
		GoroutineID: gs.GoroutineID,
		Clock:       gs.Clock.Copy(),
		StackTrace:  stacks.NewIntStackWithMap(*gs.StackTrace.GetItems().Copy(), nil),
	}
}
