package domain

import "StaticRaceDetector/utils"

type GoroutineState struct {
	GoroutineID int
	Clock       VectorClock
	StackTrace  *utils.Stack
}

type GoroutineStateJSON struct {
	GoroutineID int
	Clock       VectorClock
}

var GoroutineCounter = utils.NewCounter()

func NewEmptyGoroutineState() *GoroutineState {
	return &GoroutineState{
		Clock:       VectorClock{},
		GoroutineID: GoroutineCounter.GetNext(),
		StackTrace:  utils.NewStack(),
	}
}

func NewGoroutineExecutionState(state *GoroutineState) *GoroutineState {
	state.Increment()
	return &GoroutineState{
		Clock:       state.Clock,
		GoroutineID: GoroutineCounter.GetNext(),
		StackTrace:  state.StackTrace,
	}
}

func (gs *GoroutineState) Increment() {
	gs.Clock[gs.GoroutineID] += 1
}

func (gs *GoroutineState) MayConcurrent(state *GoroutineState) bool {
	timestampAidA, _ := gs.Clock[gs.GoroutineID]
	timestampAidB, _ := state.Clock[gs.GoroutineID]
	timestampBidA, _ := gs.Clock[state.GoroutineID]
	timestampBidB, _ := state.Clock[state.GoroutineID]
	isBefore := timestampAidA <= timestampAidB && timestampBidA < timestampBidB
	isAfter := timestampBidB <= timestampBidA && timestampAidB < timestampAidA
	return !(isBefore || isAfter)
}
func (gs *GoroutineState) Copy() *GoroutineState {
	return &GoroutineState{GoroutineID: gs.GoroutineID, Clock: gs.Clock.Copy(), StackTrace: gs.StackTrace}
}

func (gs *GoroutineState) ToJSON() *GoroutineStateJSON {
	dumpJson := GoroutineStateJSON{}
	dumpJson.GoroutineID = gs.GoroutineID
	dumpJson.Clock = gs.Clock
	return &dumpJson
}
