package ssaUtils

import (
	"StaticRaceDetector/domain"
	"golang.org/x/tools/go/ssa"
)

func HandleBuiltin(guardedAccesses *[]*domain.GuardedAccess, GoroutineState *domain.GoroutineState, callCommon *ssa.Builtin, args []ssa.Value) {

	switch name := callCommon.Name(); name {
	case "delete":
		domain.AddGuardedAccess(guardedAccesses, callCommon.Pos(), args[0], domain.GuardAccessWrite, GoroutineState)
	case "len":
	case "cap":
		domain.AddGuardedAccess(guardedAccesses, callCommon.Pos(), args[0], domain.GuardAccessRead, GoroutineState)
	case "append":
		domain.AddGuardedAccess(guardedAccesses, callCommon.Pos(), args[1], domain.GuardAccessRead, GoroutineState)
		domain.AddGuardedAccess(guardedAccesses, callCommon.Pos(), args[0], domain.GuardAccessWrite, GoroutineState)
	case "copy":
		domain.AddGuardedAccess(guardedAccesses, callCommon.Pos(), args[0], domain.GuardAccessRead, GoroutineState)
		domain.AddGuardedAccess(guardedAccesses, callCommon.Pos(), args[1], domain.GuardAccessWrite, GoroutineState)
	}
}
