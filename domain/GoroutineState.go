package domain

import "StaticRaceDetector/utils"

type GoroutineState struct {
	GoroutineID int
	Clock       VectorClock
	Lockset     *Lockset
}

type GoroutineStateJSON struct {
	GoroutineID int
	Clock       VectorClock
	LocksetJson *LocksetJson
}

var GoroutineCounter = utils.NewCounter()

func NewEmptyGoroutineState() *GoroutineState {
	return &GoroutineState{
		Clock:       VectorClock{},
		Lockset:     NewEmptyLockSet(),
		GoroutineID: GoroutineCounter.GetNext(),
	}
}

func NewGoroutineExecutionState(state *GoroutineState) *GoroutineState {
	state.Increment()
	return &GoroutineState{
		Clock:       state.Clock,
		Lockset:     state.Lockset,
		GoroutineID: GoroutineCounter.GetNext(),
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
	return &GoroutineState{GoroutineID: gs.GoroutineID, Clock: gs.Clock.Copy(), Lockset: gs.Lockset.Copy()}
}

func (gs *GoroutineState) ToJSON() *GoroutineStateJSON {
	dumpJson := GoroutineStateJSON{}
	dumpJson.LocksetJson = gs.Lockset.ToJSON()
	dumpJson.GoroutineID = gs.GoroutineID
	dumpJson.Clock = gs.Clock
	return &dumpJson
}

func (gs *GoroutineState) MergeStates(stateToMerge *GoroutineState, isConditional bool) {
	gs.Clock = stateToMerge.Clock
	if isConditional {
		gs.Lockset.UpdateLockSet(nil, stateToMerge.Lockset.ExistingUnlocks)
	} else {
		gs.Lockset.UpdateLockSet(stateToMerge.Lockset.ExistingLocks, stateToMerge.Lockset.ExistingUnlocks)
	}
}
