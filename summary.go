package main

import (
	"StaticRaceDetector/domain"
	"StaticRaceDetector/utils"
	"go/token"
	"golang.org/x/tools/go/ssa"
	"strconv"
)

var pkgNamesToCheck = []string{"pkg", "main", "StaticRaceDetector/testutils/NestedFunctions", "StaticRaceDetector/testutils/DataRaceShadowedErr", "StaticRaceDetector/testutils/Lock", "StaticRaceDetector/testutils/LockAndUnlock", "StaticRaceDetector/testutils/LockAndUnlockIfBranch", "StaticRaceDetector/testutils/DeferredLockAndUnlockIfBranch", "StaticRaceDetector/testutils/LockAndUnlockIfMap", "StaticRaceDetector/testutils/DataRaceMap", "StaticRaceDetector/testutils/DataRaceShadowedErr", "StaticRaceDetector/testutils/DataRaceProperty"}
var GuardedAccessCounter = utils.NewCounter()

func addGuardedAccess(guardedAccesses *[]*domain.GuardedAccess, pos token.Pos, value ssa.Value, kind domain.OpKind, GoroutineState *domain.GoroutineState) {
	GoroutineState.Increment()
	guardedAccessToAdd := &domain.GuardedAccess{ID: GuardedAccessCounter.GetNext(), Pos: pos, Value: value, OpKind: kind, State: GoroutineState.Copy()}
	*guardedAccesses = append(*guardedAccesses, guardedAccessToAdd)
}
func GetBlockSummary(block *ssa.BasicBlock, GoroutineState *domain.GoroutineState) ([]*ssa.CallCommon, []*domain.GuardedAccess, *domain.GoroutineState) {
	deferredFunctions := make([]*ssa.CallCommon, 0)
	guardedAccesses := make([]*domain.GuardedAccess, 0)
	instrs := utils.FilterDebug(block.Instrs)
	for _, ins := range instrs {
		switch call := ins.(type) {
		case *ssa.UnOp:
			addGuardedAccess(&guardedAccesses, call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		case *ssa.Field:
			addGuardedAccess(&guardedAccesses, call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		case *ssa.FieldAddr:
			addGuardedAccess(&guardedAccesses, call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		case *ssa.Index:
			addGuardedAccess(&guardedAccesses, call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		case *ssa.IndexAddr:
			addGuardedAccess(&guardedAccesses, call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		case *ssa.Lookup:
			addGuardedAccess(&guardedAccesses, call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		case *ssa.Panic:
			addGuardedAccess(&guardedAccesses, call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		case *ssa.Range:
			addGuardedAccess(&guardedAccesses, call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		case *ssa.TypeAssert:
			addGuardedAccess(&guardedAccesses, call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		case *ssa.BinOp:
			addGuardedAccess(&guardedAccesses, call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
			addGuardedAccess(&guardedAccesses, call.Pos(), call.Y, domain.GuardAccessRead, GoroutineState)
		case *ssa.If:
			addGuardedAccess(&guardedAccesses, call.Pos(), call.Cond, domain.GuardAccessRead, GoroutineState)
		case *ssa.MapUpdate:
			addGuardedAccess(&guardedAccesses, call.Pos(), call.Map, domain.GuardAccessWrite, GoroutineState)
			addGuardedAccess(&guardedAccesses, call.Pos(), call.Value, domain.GuardAccessRead, GoroutineState)
		case *ssa.Store:
			addGuardedAccess(&guardedAccesses, call.Pos(), call.Val, domain.GuardAccessRead, GoroutineState)
			addGuardedAccess(&guardedAccesses, call.Pos(), call.Addr, domain.GuardAccessWrite, GoroutineState)
		case *ssa.Call:
			callCommon := call.Common()
			guardedAccessesRet, guardedState := GetFunctionSummary(callCommon, GoroutineState.Copy())
			guardedAccesses = append(guardedAccesses, guardedAccessesRet...)
			GoroutineState.MergeStates(guardedState, false)
		case *ssa.Go:
			callCommon := call.Common()
			GoroutineState.Increment()
			newState := domain.NewGoroutineState()
			newState.Clock = GoroutineState.Clock
			guardedAccessesRet, _ := GetFunctionSummary(callCommon, newState.Copy())
			guardedAccesses = append(guardedAccesses, guardedAccessesRet...)
		case *ssa.Defer:
			callCommon := call.Common()
			deferredFunctions = append(deferredFunctions, callCommon)
		}
	}
	return deferredFunctions, guardedAccesses, GoroutineState
}

func GetFunctionSummary(callCommon *ssa.CallCommon, GoroutineState *domain.GoroutineState) ([]*domain.GuardedAccess, *domain.GoroutineState) {
	var conditionalBlocks = map[string]struct{}{
		"if.then":     {},
		"if.else":     {},
		"select.body": {},
	}
	guardedAccesses := make([]*domain.GuardedAccess, 0)

	if utils.IsCallToAny(callCommon, "(*sync.Mutex).Lock") {
		receiver := callCommon.Args[0]
		LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
		locks := map[string]*ssa.CallCommon{LockName: callCommon}
		GoroutineState.Lockset.UpdateLockSet(locks, nil)
		return guardedAccesses, GoroutineState
	}
	if utils.IsCallToAny(callCommon, "(*sync.Mutex).Unlock") {
		receiver := callCommon.Args[0]
		LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
		locks := map[string]*ssa.CallCommon{LockName: callCommon}
		GoroutineState.Lockset.UpdateLockSet(nil, locks)
		return guardedAccesses, GoroutineState
	}
	if utils.IsCallToAny(callCommon, "delete") {
		addGuardedAccess(&guardedAccesses, callCommon.Pos(), callCommon.Args[0], domain.GuardAccessWrite, GoroutineState)
		return guardedAccesses, GoroutineState
	}
	if utils.IsCallToAny(callCommon, "len") || utils.IsCallToAny(callCommon, "cap") {
		addGuardedAccess(&guardedAccesses, callCommon.Pos(), callCommon.Args[0], domain.GuardAccessRead, GoroutineState)
		return guardedAccesses, GoroutineState
	}
	if utils.IsCallToAny(callCommon, "append") {
		addGuardedAccess(&guardedAccesses, callCommon.Pos(), callCommon.Args[1], domain.GuardAccessRead, GoroutineState)
		addGuardedAccess(&guardedAccesses, callCommon.Pos(), callCommon.Args[0], domain.GuardAccessWrite, GoroutineState)
		return guardedAccesses, GoroutineState
	}
	if utils.IsCallToAny(callCommon, "copy") {
		addGuardedAccess(&guardedAccesses, callCommon.Pos(), callCommon.Args[0], domain.GuardAccessRead, GoroutineState)
		addGuardedAccess(&guardedAccesses, callCommon.Pos(), callCommon.Args[1], domain.GuardAccessWrite, GoroutineState)
		return guardedAccesses, GoroutineState
	}

	fn, ok := callCommon.Value.(*ssa.Function)
	if !ok {
		fn = callCommon.Value.(*ssa.MakeClosure).Fn.(*ssa.Function)
	}

	pkgName := fn.Pkg.Pkg.Path()
	found := false
	for _, pkgNameToCheck := range pkgNamesToCheck {
		if pkgName == pkgNameToCheck {
			found = true
			break
		}
	}
	if !found {
		return guardedAccesses, GoroutineState
	}

	deferredFunctions := make([]*domain.ConditionalFunction, 0)
	for _, block := range fn.Blocks {
		deferredFunctionsRet, guardedAccessesRet, goroutineState := GetBlockSummary(block, GoroutineState.Copy())
		guardedAccesses = append(guardedAccesses, guardedAccessesRet...)

		if _, ok := conditionalBlocks[block.Comment]; ok {
			GoroutineState.MergeStates(goroutineState, true) // Ignore locks in a condition branch since it's a must set.
			for _, deferredFunctionRet := range deferredFunctionsRet {
				deferredFunctions = append(deferredFunctions, &domain.ConditionalFunction{IsConditional: true, Function: deferredFunctionRet})
			}
		} else {
			GoroutineState.MergeStates(goroutineState, false)
			for _, deferredFunctionRet := range deferredFunctionsRet {
				deferredFunctions = append(deferredFunctions, &domain.ConditionalFunction{IsConditional: false, Function: deferredFunctionRet})
			}
		}
	}

	for i := len(deferredFunctions) - 1; i >= 0; i-- {
		res, GoroutineStateRet := GetFunctionSummary(deferredFunctions[i].Function, GoroutineState.Copy())
		guardedAccesses = append(guardedAccesses, res...)
		GoroutineState.MergeStates(GoroutineStateRet, deferredFunctions[i].IsConditional)
	}
	return guardedAccesses, GoroutineState
}
