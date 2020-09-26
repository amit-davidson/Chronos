package ssaUtils

import (
	"StaticRaceDetector/domain"
	"golang.org/x/tools/go/ssa"
	"strconv"
)

func AddLock(GoroutineState *domain.GoroutineState, call *ssa.CallCommon, isUnlock bool) *domain.GoroutineState {
	receiver := call.Args[0]
	LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
	lock := map[string]*ssa.CallCommon{LockName: call}
	if isUnlock {
		GoroutineState.Lockset.UpdateLockSet(nil, lock)
	} else {
		GoroutineState.Lockset.UpdateLockSet(lock, nil)
	}
	return GoroutineState
}