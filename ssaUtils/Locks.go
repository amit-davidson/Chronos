package ssaUtils

import (
	"StaticRaceDetector/domain"
	"golang.org/x/tools/go/ssa"
	"strconv"
)

func AddLock(funcState *domain.FunctionState, call *ssa.CallCommon, isUnlock bool) {
	receiver := call.Args[0]
	LockName := strconv.Itoa(int(receiver.Pos()))
	lock := map[string]*ssa.CallCommon{LockName: call}
	if isUnlock {
		funcState.Lockset.UpdateLockSet(nil, lock)
	} else {
		funcState.Lockset.UpdateLockSet(lock, nil)
	}
}