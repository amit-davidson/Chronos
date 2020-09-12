package main

import (
	"StaticRaceDetector/utils"
	"golang.org/x/tools/go/ssa"
	"strconv"
	"sync"
)

var pkgNamesToCheck = []string{"pkg", "main"}

var wg sync.WaitGroup

func addGuardedAccess(guardedAccesses *[]*guardedAccess, value ssa.Value, kind opKind, currentLockset *lockset, goroutineId string) {
	guardedAccessToAdd := &guardedAccess{id: utils.GetUUID(), value: value, opKind: kind, lockset: currentLockset.Copy(), GoroutineId: goroutineId}
	*guardedAccesses = append(*guardedAccesses, guardedAccessToAdd)
}
func GetBlockSummary(block *ssa.BasicBlock, ls *lockset, goroutineId string) (*lockset, []*lockset, []*guardedAccess) {
	deferredCalls := make([]*lockset, 0)
	guardedAccesses := make([]*guardedAccess, 0)
	instrs := FilterDebug(block.Instrs)
	for _, ins := range instrs {
		switch call := ins.(type) {
		case *ssa.UnOp:
			addGuardedAccess(&guardedAccesses, call.X, read, ls, goroutineId)
		case *ssa.Field:
			addGuardedAccess(&guardedAccesses, call.X, read, ls, goroutineId)
		case *ssa.FieldAddr:
			addGuardedAccess(&guardedAccesses, call.X, read, ls, goroutineId)
		case *ssa.Index:
			addGuardedAccess(&guardedAccesses, call.X, read, ls, goroutineId)
		case *ssa.IndexAddr:
			addGuardedAccess(&guardedAccesses, call.X, read, ls, goroutineId)
		case *ssa.Lookup:
			addGuardedAccess(&guardedAccesses, call.X, read, ls, goroutineId)
		case *ssa.Panic:
			addGuardedAccess(&guardedAccesses, call.X, read, ls, goroutineId)
		case *ssa.Range:
			addGuardedAccess(&guardedAccesses, call.X, read, ls, goroutineId)
		case *ssa.TypeAssert:
			addGuardedAccess(&guardedAccesses, call.X, read, ls, goroutineId)
		//case *ssa.Send:
		//	addGuardedAccess(&guardedAccesses, call.X, read, ls, goroutineId)
		case *ssa.BinOp:
			addGuardedAccess(&guardedAccesses, call.X, read, ls, goroutineId)
			addGuardedAccess(&guardedAccesses, call.Y, read, ls, goroutineId)
		case *ssa.If:
			addGuardedAccess(&guardedAccesses, call.Cond, read, ls, goroutineId)
		case *ssa.MapUpdate:
			addGuardedAccess(&guardedAccesses, call.Map, write, ls, goroutineId)
			addGuardedAccess(&guardedAccesses, call.Value, read, ls, goroutineId)
		case *ssa.Store:
			addGuardedAccess(&guardedAccesses, call.Addr, write, ls, goroutineId)
			addGuardedAccess(&guardedAccesses, call.Val, read, ls, goroutineId)
		case *ssa.Call:
			callCommon := call.Common()
			if IsCallToAny(callCommon, "(*sync.Mutex).Lock") {
				receiver := callCommon.Args[0]
				LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
				locks := map[string]*ssa.CallCommon{LockName: callCommon}
				ls.updateLockSet(locks, nil)
			}
			if IsCallToAny(callCommon, "(*sync.Mutex).Unlock") {
				receiver := callCommon.Args[0]
				LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
				locks := map[string]*ssa.CallCommon{LockName: callCommon}
				ls.updateLockSet(nil, locks)
			}
			if IsCallToAny(callCommon, "delete") {
				addGuardedAccess(&guardedAccesses, callCommon.Args[0], write, ls, goroutineId)
			}
			if IsCallToAny(callCommon, "len") || IsCallToAny(callCommon, "cap") {
				addGuardedAccess(&guardedAccesses, callCommon.Args[0], read, ls, goroutineId)
			}
			if IsCallToAny(callCommon, "append") {
				addGuardedAccess(&guardedAccesses, callCommon.Args[0], write, ls, goroutineId)
				addGuardedAccess(&guardedAccesses, callCommon.Args[1], read, ls, goroutineId)
			}
			if IsCallToAny(callCommon, "copy") {
				addGuardedAccess(&guardedAccesses, callCommon.Args[0], read, ls, goroutineId)
				addGuardedAccess(&guardedAccesses, callCommon.Args[1], write, ls, goroutineId)
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
				ls.updateLockSet(lsRet.existingLocks, lsRet.existingUnlocks)
			}
			continue
		case *ssa.Go:
			callCommon := call.Common()
			function, ok := callCommon.Value.(*ssa.Function)
			if !ok {
				function = callCommon.Value.(*ssa.MakeClosure).Fn.(*ssa.Function)
			}

			//wg.Add(1)
			lsRet, guardedAccessesRet := GetFunctionSummary(function, newEmptyLockSet(), utils.GetUUID())
			guardedAccesses = append(guardedAccesses, guardedAccessesRet...)
			ls.updateLockSet(lsRet.existingLocks, lsRet.existingUnlocks)
		case *ssa.Defer:
			callCommon := call.Common()
			if IsCallToAny(callCommon, "(*sync.Mutex).Lock") {
				receiver := callCommon.Args[0]
				LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
				locks := map[string]*ssa.CallCommon{LockName: callCommon}
				deferredCalls = append(deferredCalls, newLockSet(locks, nil))
			}
			if IsCallToAny(callCommon, "(*sync.Mutex).Unlock") {
				receiver := callCommon.Args[0]
				LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
				locks := map[string]*ssa.CallCommon{LockName: callCommon}
				deferredCalls = append(deferredCalls, newLockSet(nil, locks))
			}
			continue
		}
	}
	return ls, deferredCalls, guardedAccesses
}

func GetFunctionSummary(fn *ssa.Function, ls *lockset, goroutineIdCounter string) (*lockset, []*guardedAccess) {
	var conditionalBlocks = map[string]struct{}{
		"if.then": {},
		"if.else": {},
		"select.body": {},
	}

	deferredCalls := make([]*lockset, 0)
	guardedAccesses := make([]*guardedAccess, 0)
	for _, block := range fn.Blocks {
		lsRet, deferredCallsRet, guardedAccessesRet := GetBlockSummary(block, ls.Copy(), goroutineIdCounter)
		guardedAccesses = append(guardedAccesses, guardedAccessesRet...)
		if _, ok := conditionalBlocks[block.Comment]; ok {
			ls.updateLockSet(nil, lsRet.existingUnlocks) // Ignore locks in a condition branch since it's a must set.
			for _, deferredCallRet := range deferredCallsRet {
				deferredCalls = append(deferredCalls, newLockSet(nil, deferredCallRet.existingUnlocks))
			}
		} else {
			ls.updateLockSet(lsRet.existingLocks, lsRet.existingUnlocks)
			for _, deferredCallRet := range deferredCallsRet {
				deferredCalls = append(deferredCalls, newLockSet(deferredCallRet.existingLocks, deferredCallRet.existingUnlocks))
			}
		}
	}

	for i := len(deferredCalls) - 1; i >= 0; i-- {
		ls.updateLockSet(deferredCalls[i].existingLocks, deferredCalls[i].existingUnlocks)
	}
	//wg.Done()
	return ls, guardedAccesses
}
