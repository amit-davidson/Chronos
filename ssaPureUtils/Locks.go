package ssaPureUtils

import (
	"github.com/amit-davidson/Chronos/utils"
	"golang.org/x/tools/go/ssa"
)

func IsLock(call *ssa.Function) bool {
	return utils.IsCallTo(call, "(*sync.Mutex).Lock")
}


func IsUnlock(call *ssa.Function) bool {
	return utils.IsCallTo(call, "(*sync.Mutex).Unlock")
}
