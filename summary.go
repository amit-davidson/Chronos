package main

import (
	"StaticRaceDetector/domain"
	"StaticRaceDetector/utils"
	"go/token"
	"golang.org/x/tools/go/ssa"
	"strconv"
)

var pkgNamesToCheck = []string{"pkg", "main", "StaticRaceDetector/testutils/NestedFunctions", "StaticRaceDetector/testutils/DataRaceShadowedErr"}
var GuardedAccessCounter = utils.NewCounter()

func addGuardedAccess(guardedAccesses *[]*domain.GuardedAccess, pos token.Pos ,value ssa.Value, kind domain.OpKind, GoroutineState *domain.GoroutineState) {
	GoroutineState.Increment()
	guardedAccessToAdd := &domain.GuardedAccess{ID: GuardedAccessCounter.GetNext(), Pos: pos, Value: value, OpKind: kind, State: GoroutineState.Copy()}
	*guardedAccesses = append(*guardedAccesses, guardedAccessToAdd)
}
func GetBlockSummary(block *ssa.BasicBlock, GoroutineState *domain.GoroutineState) (*domain.GoroutineState, []*domain.Lockset, []*domain.GuardedAccess) {
	deferredCalls := make([]*domain.Lockset, 0)
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
			if utils.IsCallToAny(callCommon, "(*sync.Mutex).Lock") {
				receiver := callCommon.Args[0]
				LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
				locks := map[string]*ssa.CallCommon{LockName: callCommon}
				GoroutineState.Lockset.UpdateLockSet(locks, nil)
			}
			if utils.IsCallToAny(callCommon, "(*sync.Mutex).Unlock") {
				receiver := callCommon.Args[0]
				LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
				locks := map[string]*ssa.CallCommon{LockName: callCommon}
				GoroutineState.Lockset.UpdateLockSet(nil, locks)
			}
			if utils.IsCallToAny(callCommon, "delete") {
				addGuardedAccess(&guardedAccesses, callCommon.Pos(), callCommon.Args[0], domain.GuardAccessWrite, GoroutineState)
			}
			if utils.IsCallToAny(callCommon, "len") || utils.IsCallToAny(callCommon, "cap") {
				addGuardedAccess(&guardedAccesses, callCommon.Pos(), callCommon.Args[0], domain.GuardAccessRead, GoroutineState)
			}
			if utils.IsCallToAny(callCommon, "append") {
				addGuardedAccess(&guardedAccesses, callCommon.Pos(), callCommon.Args[1], domain.GuardAccessRead, GoroutineState)
				addGuardedAccess(&guardedAccesses, callCommon.Pos(), callCommon.Args[0], domain.GuardAccessWrite, GoroutineState)
			}
			if utils.IsCallToAny(callCommon, "copy") {
				addGuardedAccess(&guardedAccesses, callCommon.Pos(), callCommon.Args[0], domain.GuardAccessRead, GoroutineState)
				addGuardedAccess(&guardedAccesses, callCommon.Pos(), callCommon.Args[1], domain.GuardAccessWrite, GoroutineState)
			}
			if function, isFunctionCall := callCommon.Value.(*ssa.Function); isFunctionCall {
				pkgName := function.Pkg.Pkg.Path()
				found := false
				for _, pkgNameToCheck := range pkgNamesToCheck {
					if pkgName == pkgNameToCheck {
						found = true
						break
					}
				}
				if !found {
					continue
				}
				lsRet, guardedAccessesRet := GetFunctionSummary(function, GoroutineState)
				guardedAccesses = append(guardedAccesses, guardedAccessesRet...)
				GoroutineState.Lockset.UpdateLockSet(lsRet.ExistingLocks, lsRet.ExistingUnlocks)
			}
			continue
		case *ssa.Go:
			callCommon := call.Common()
			function, ok := callCommon.Value.(*ssa.Function)
			if !ok {
				function = callCommon.Value.(*ssa.MakeClosure).Fn.(*ssa.Function)
			}

			newState := domain.NewGoroutineState()
			GoroutineState.Increment()
			newState.Clock = GoroutineState.Clock
			lsRet, guardedAccessesRet := GetFunctionSummary(function, newState.Copy())
			guardedAccesses = append(guardedAccesses, guardedAccessesRet...)
			GoroutineState.Lockset.UpdateLockSet(lsRet.ExistingLocks, lsRet.ExistingUnlocks)
		case *ssa.Defer:
			callCommon := call.Common()
			if utils.IsCallToAny(callCommon, "(*sync.Mutex).Lock") {
				receiver := callCommon.Args[0]
				LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
				locks := map[string]*ssa.CallCommon{LockName: callCommon}
				deferredCalls = append(deferredCalls, domain.NewLockSet(locks, nil))
			}
			if utils.IsCallToAny(callCommon, "(*sync.Mutex).Unlock") {
				receiver := callCommon.Args[0]
				LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
				locks := map[string]*ssa.CallCommon{LockName: callCommon}
				deferredCalls = append(deferredCalls, domain.NewLockSet(nil, locks))
			}
			continue
		}
	}
	return GoroutineState, deferredCalls, guardedAccesses
}

func GetFunctionSummary(fn *ssa.Function, GoroutineState *domain.GoroutineState) (*domain.Lockset, []*domain.GuardedAccess) {
	var conditionalBlocks = map[string]struct{}{
		"if.then":     {},
		"if.else":     {},
		"select.body": {},
	}

	deferredCalls := make([]*domain.Lockset, 0)
	guardedAccesses := make([]*domain.GuardedAccess, 0)
	for _, block := range fn.Blocks {
		// We copy the lockset since changes to it aren't determined outside
		updatedGoroutineState, deferredCallsRet, guardedAccessesRet := GetBlockSummary(block, GoroutineState.Copy())
		guardedAccesses = append(guardedAccesses, guardedAccessesRet...)
		GoroutineState.Clock = updatedGoroutineState.Clock
		if _, ok := conditionalBlocks[block.Comment]; ok {
			GoroutineState.Lockset.UpdateLockSet(nil, updatedGoroutineState.Lockset.ExistingUnlocks) // Ignore locks in a condition branch since it's a must set.
			for _, deferredCallRet := range deferredCallsRet {
				deferredCalls = append(deferredCalls, domain.NewLockSet(nil, deferredCallRet.ExistingUnlocks))
			}
		} else {
			GoroutineState.Lockset.UpdateLockSet(updatedGoroutineState.Lockset.ExistingLocks, updatedGoroutineState.Lockset.ExistingUnlocks)
			for _, deferredCallRet := range deferredCallsRet {
				deferredCalls = append(deferredCalls, domain.NewLockSet(deferredCallRet.ExistingLocks, deferredCallRet.ExistingUnlocks))
			}
		}
	}

	for i := len(deferredCalls) - 1; i >= 0; i-- {
		GoroutineState.Lockset.UpdateLockSet(deferredCalls[i].ExistingLocks, deferredCalls[i].ExistingUnlocks)
	}
	return GoroutineState.Lockset, guardedAccesses
}
