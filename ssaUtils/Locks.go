package ssaUtils

import (
	"github.com/amit-davidson/Chronos/domain"
	"go/token"
	"golang.org/x/tools/go/ssa"
)

func AddLock(funcState *domain.BlockState, call *ssa.CallCommon, isUnlock bool) {
	receiver := call.Args[0]
	LockName := receiver.Pos()
	lock := map[token.Pos]*ssa.CallCommon{LockName: call}
	if isUnlock {
		funcState.Lockset.UpdateLockSet(nil, lock)
	} else {
		funcState.Lockset.UpdateLockSet(lock, nil)
	}
}
