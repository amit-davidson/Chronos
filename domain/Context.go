package domain

import "github.com/amit-davidson/Chronos/utils"

type Context struct {
	GoroutineID int
	Clock       VectorClock
	StackTrace  *utils.Stack
}

type ContextJSON struct {
	GoroutineID int
	Clock       VectorClock
}

var GoroutineCounter = utils.NewCounter()

func NewEmptyContext() *Context {
	return &Context{
		Clock:       VectorClock{},
		GoroutineID: GoroutineCounter.GetNext(),
		StackTrace:  utils.NewStack(),
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
	timestampAidA, _ := gs.Clock[gs.GoroutineID]
	timestampAidB, _ := state.Clock[gs.GoroutineID]
	timestampBidA, _ := gs.Clock[state.GoroutineID]
	timestampBidB, _ := state.Clock[state.GoroutineID]
	isBefore := timestampAidA <= timestampAidB && timestampBidA < timestampBidB
	isAfter := timestampBidB <= timestampBidA && timestampAidB < timestampAidA
	return !(isBefore || isAfter)
}
func (gs *Context) Copy() *Context {
	return &Context{GoroutineID: gs.GoroutineID, Clock: gs.Clock.Copy(), StackTrace: gs.StackTrace}
}

func (gs *Context) ToJSON() *ContextJSON {
	dumpJson := ContextJSON{}
	dumpJson.GoroutineID = gs.GoroutineID
	dumpJson.Clock = gs.Clock
	return &dumpJson
}
