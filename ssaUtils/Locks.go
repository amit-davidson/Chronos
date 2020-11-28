package ssaUtils

import (
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/ssaPureUtils"
	"go/token"
	"golang.org/x/tools/go/ssa"
)

func AddLock(funcState *domain.BlockState, call *ssa.CallCommon, isUnlock bool) {
	recv := call.Args[0]
	mutexPos := ssaPureUtils.GetMutexPos(recv)
	lock := map[token.Pos]*ssa.CallCommon{mutexPos: call}
	if isUnlock {
		funcState.Lockset.UpdateWithNewLockSet(nil, lock)
	} else {
		funcState.Lockset.UpdateWithNewLockSet(lock, nil)
	}
}