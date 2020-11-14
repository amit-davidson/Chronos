package ssaUtils

import (
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/ssaPureUtils"
	"github.com/amit-davidson/Chronos/utils"
	"go/token"
	"golang.org/x/tools/go/ssa"
)

func AddLock(funcState *domain.BlockState, call *ssa.CallCommon, isUnlock bool) {
	recv := call.Args[0]
	mutexPos := ssaPureUtils.GetMutexPos(recv)
	lock := map[token.Pos]*ssa.CallCommon{mutexPos: call}
	if isUnlock {
		funcState.Lockset.UpdateLockSet(nil, lock)
	} else {
		funcState.Lockset.UpdateLockSet(lock, nil)
	}
}

func IsLock(call *ssa.Function) bool {
	return utils.IsCallTo(call, "(*sync.Mutex).Lock")
}


func IsUnlock(call *ssa.Function) bool {
	return utils.IsCallTo(call, "(*sync.Mutex).Unlock")
}
