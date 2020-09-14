package main

import (
	"StaticRaceDetector/domain"
	"StaticRaceDetector/utils"
	"golang.org/x/tools/go/ssa"
	"strconv"
)

var pkgNamesToCheck = []string{"pkg", "main"}

func addGuardedAccess(guardedAccesses *[]*domain.GuardedAccess, value ssa.Value, kind domain.OpKind, currentLockset *domain.Lockset, goroutineId string) {
	guardedAccessToAdd := &domain.GuardedAccess{ID: utils.GetUUID(), Value: value, OpKind: kind, Lockset: currentLockset.Copy(), GoroutineId: goroutineId}
	*guardedAccesses = append(*guardedAccesses, guardedAccessToAdd)
}
func GetBlockSummary(block *ssa.BasicBlock, ls *domain.Lockset, goroutineId string) (*domain.Lockset, []*domain.Lockset, []*domain.GuardedAccess) {
	deferredCalls := make([]*domain.Lockset, 0)
	guardedAccesses := make([]*domain.GuardedAccess, 0)
	instrs := FilterDebug(block.Instrs)
	for _, ins := range instrs {
		switch call := ins.(type) {
		case *ssa.UnOp:
			addGuardedAccess(&guardedAccesses, call.X, domain.GuardAccessRead, ls, goroutineId)
		case *ssa.Field:
			addGuardedAccess(&guardedAccesses, call.X, domain.GuardAccessRead, ls, goroutineId)
		case *ssa.FieldAddr:
			addGuardedAccess(&guardedAccesses, call.X, domain.GuardAccessRead, ls, goroutineId)
		case *ssa.Index:
			addGuardedAccess(&guardedAccesses, call.X, domain.GuardAccessRead, ls, goroutineId)
		case *ssa.IndexAddr:
			addGuardedAccess(&guardedAccesses, call.X, domain.GuardAccessRead, ls, goroutineId)
		case *ssa.Lookup:
			addGuardedAccess(&guardedAccesses, call.X, domain.GuardAccessRead, ls, goroutineId)
		case *ssa.Panic:
			addGuardedAccess(&guardedAccesses, call.X, domain.GuardAccessRead, ls, goroutineId)
		case *ssa.Range:
			addGuardedAccess(&guardedAccesses, call.X, domain.GuardAccessRead, ls, goroutineId)
		case *ssa.TypeAssert:
			addGuardedAccess(&guardedAccesses, call.X, domain.GuardAccessRead, ls, goroutineId)
		//case *ssa.Send:
		//	addGuardedAccess(&guardedAccesses, call.X, read, ls, goroutineId)
		case *ssa.BinOp:
			addGuardedAccess(&guardedAccesses, call.X, domain.GuardAccessRead, ls, goroutineId)
			addGuardedAccess(&guardedAccesses, call.Y, domain.GuardAccessRead, ls, goroutineId)
		case *ssa.If:
			addGuardedAccess(&guardedAccesses, call.Cond, domain.GuardAccessRead, ls, goroutineId)
		case *ssa.MapUpdate:
			addGuardedAccess(&guardedAccesses, call.Map, domain.GuardAccessRead, ls, goroutineId)
			addGuardedAccess(&guardedAccesses, call.Value, domain.GuardAccessRead, ls, goroutineId)
		case *ssa.Store:
			addGuardedAccess(&guardedAccesses, call.Addr, domain.GuardAccessWrite, ls, goroutineId)
			addGuardedAccess(&guardedAccesses, call.Val, domain.GuardAccessRead, ls, goroutineId)
		case *ssa.Call:
			callCommon := call.Common()
			if IsCallToAny(callCommon, "(*sync.Mutex).Lock") {
				receiver := callCommon.Args[0]
				LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
				locks := map[string]*ssa.CallCommon{LockName: callCommon}
				ls.UpdateLockSet(locks, nil)
			}
			if IsCallToAny(callCommon, "(*sync.Mutex).Unlock") {
				receiver := callCommon.Args[0]
				LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
				locks := map[string]*ssa.CallCommon{LockName: callCommon}
				ls.UpdateLockSet(nil, locks)
			}
			if IsCallToAny(callCommon, "delete") {
				addGuardedAccess(&guardedAccesses, callCommon.Args[0], domain.GuardAccessWrite, ls, goroutineId)
			}
			if IsCallToAny(callCommon, "len") || IsCallToAny(callCommon, "cap") {
				addGuardedAccess(&guardedAccesses, callCommon.Args[0], domain.GuardAccessRead, ls, goroutineId)
			}
			if IsCallToAny(callCommon, "append") {
				addGuardedAccess(&guardedAccesses, callCommon.Args[0], domain.GuardAccessWrite, ls, goroutineId)
				addGuardedAccess(&guardedAccesses, callCommon.Args[1], domain.GuardAccessRead, ls, goroutineId)
			}
			if IsCallToAny(callCommon, "copy") {
				addGuardedAccess(&guardedAccesses, callCommon.Args[0], domain.GuardAccessRead, ls, goroutineId)
				addGuardedAccess(&guardedAccesses, callCommon.Args[1], domain.GuardAccessWrite, ls, goroutineId)
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
				lsRet, guardedAccessesRet := GetFunctionSummary(function, ls, goroutineId)
				guardedAccesses = append(guardedAccesses, guardedAccessesRet...)
				ls.UpdateLockSet(lsRet.ExistingLocks, lsRet.ExistingUnlocks)
			}
			continue
		case *ssa.Go:
			callCommon := call.Common()
			function, ok := callCommon.Value.(*ssa.Function)
			if !ok {
				function = callCommon.Value.(*ssa.MakeClosure).Fn.(*ssa.Function)
			}

			lsRet, guardedAccessesRet := GetFunctionSummary(function, domain.NewEmptyLockSet(), utils.GetUUID())
			guardedAccesses = append(guardedAccesses, guardedAccessesRet...)
			ls.UpdateLockSet(lsRet.ExistingLocks, lsRet.ExistingUnlocks)
		case *ssa.Defer:
			callCommon := call.Common()
			if IsCallToAny(callCommon, "(*sync.Mutex).Lock") {
				receiver := callCommon.Args[0]
				LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
				locks := map[string]*ssa.CallCommon{LockName: callCommon}
				deferredCalls = append(deferredCalls, domain.NewLockSet(locks, nil))
			}
			if IsCallToAny(callCommon, "(*sync.Mutex).Unlock") {
				receiver := callCommon.Args[0]
				LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
				locks := map[string]*ssa.CallCommon{LockName: callCommon}
				deferredCalls = append(deferredCalls, domain.NewLockSet(nil, locks))
			}
			continue
		}
	}
	return ls, deferredCalls, guardedAccesses
}

func GetFunctionSummary(fn *ssa.Function, ls *domain.Lockset, goroutineIdCounter string) (*domain.Lockset, []*domain.GuardedAccess) {
	var conditionalBlocks = map[string]struct{}{
		"if.then":     {},
		"if.else":     {},
		"select.body": {},
	}

	deferredCalls := make([]*domain.Lockset, 0)
	guardedAccesses := make([]*domain.GuardedAccess, 0)
	for _, block := range fn.Blocks {
		lsRet, deferredCallsRet, guardedAccessesRet := GetBlockSummary(block, ls.Copy(), goroutineIdCounter)
		guardedAccesses = append(guardedAccesses, guardedAccessesRet...)
		if _, ok := conditionalBlocks[block.Comment]; ok {
			ls.UpdateLockSet(nil, lsRet.ExistingUnlocks) // Ignore locks in a condition branch since it's a must set.
			for _, deferredCallRet := range deferredCallsRet {
				deferredCalls = append(deferredCalls, domain.NewLockSet(nil, deferredCallRet.ExistingUnlocks))
			}
		} else {
			ls.UpdateLockSet(lsRet.ExistingLocks, lsRet.ExistingUnlocks)
			for _, deferredCallRet := range deferredCallsRet {
				deferredCalls = append(deferredCalls, domain.NewLockSet(deferredCallRet.ExistingLocks, deferredCallRet.ExistingUnlocks))
			}
		}
	}

	for i := len(deferredCalls) - 1; i >= 0; i-- {
		ls.UpdateLockSet(deferredCalls[i].ExistingLocks, deferredCalls[i].ExistingUnlocks)
	}
	return ls, guardedAccesses
}
