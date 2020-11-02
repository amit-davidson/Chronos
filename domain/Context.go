package domain

import (
	"github.com/amit-davidson/Chronos/utils"
	"github.com/amit-davidson/Chronos/utils/stacks"
)

type Context struct {
	GoroutineID int
	Clock       VectorClock
	StackTrace  *stacks.IntStack

	GoroutineCounter     *utils.Counter
	GuardedAccessCounter *utils.Counter
	PosIDCounter         *utils.Counter
}

func NewEmptyContext() *Context {
	var GoroutineCounter = utils.NewCounter()
	return &Context{
		Clock:                VectorClock{},
		GoroutineID:          GoroutineCounter.GetNext(),
		StackTrace:           stacks.NewIntStack(),
		GoroutineCounter:     GoroutineCounter,
		GuardedAccessCounter: utils.NewCounter(),
		PosIDCounter:         utils.NewCounter(),
	}
}

func NewGoroutineExecutionState(state *Context) *Context {
	state.Increment()
	return &Context{
		Clock:                state.Clock,
		GoroutineID:          state.GoroutineCounter.GetNext(),
		StackTrace:           state.StackTrace,
		GoroutineCounter:     state.GoroutineCounter,
		GuardedAccessCounter: state.GuardedAccessCounter,
		PosIDCounter:         state.PosIDCounter,
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
		GoroutineID:          gs.GoroutineID,
		Clock:                gs.Clock.Copy(),
		StackTrace:           gs.StackTrace.Copy(),
		GuardedAccessCounter: gs.GuardedAccessCounter,
		PosIDCounter:         gs.PosIDCounter,
		GoroutineCounter:     gs.GoroutineCounter,
	}
}
