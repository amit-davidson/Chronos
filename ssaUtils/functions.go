package ssaUtils

import (
	"StaticRaceDetector/domain"
	"StaticRaceDetector/utils"
	"golang.org/x/tools/go/ssa"
	"strconv"
	"strings"
)

func HandleBuiltin(GoroutineState *domain.GoroutineState, callCommon *ssa.Builtin, args []ssa.Value) ([]*domain.GuardedAccess, *domain.GoroutineState) {
	guardedAccesses := make([]*domain.GuardedAccess, 0)
	switch name := callCommon.Name(); name {
	case "delete":
		return append(guardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[0], domain.GuardAccessWrite, GoroutineState)), GoroutineState
	case "cap", "len":
		return append(guardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[0], domain.GuardAccessRead, GoroutineState)), GoroutineState
	case "append":
		guardedAccesses = append(guardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[1], domain.GuardAccessRead, GoroutineState))
		return append(guardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[0], domain.GuardAccessWrite, GoroutineState)), GoroutineState
	case "copy":
		guardedAccesses = append(guardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[0], domain.GuardAccessRead, GoroutineState))
		return append(guardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[1], domain.GuardAccessWrite, GoroutineState)), GoroutineState
	}
	return guardedAccesses, GoroutineState
}

func GetBlockSummary(block *ssa.BasicBlock, GoroutineState *domain.GoroutineState) ([]*ssa.CallCommon, []*domain.GuardedAccess, *domain.GoroutineState) {
	deferredFunctions := make([]*ssa.CallCommon, 0)
	guardedAccesses := make([]*domain.GuardedAccess, 0)
	for _, ins := range block.Instrs {
		switch call := ins.(type) {
		case *ssa.UnOp:
			guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
			guardedAccesses = append(guardedAccesses, guardedAccess)
		case *ssa.Field:
			guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
			guardedAccesses = append(guardedAccesses, guardedAccess)
		case *ssa.FieldAddr:
			guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
			guardedAccesses = append(guardedAccesses, guardedAccess)
		case *ssa.Index:
			guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
			guardedAccesses = append(guardedAccesses, guardedAccess)
		case *ssa.IndexAddr:
			guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
			guardedAccesses = append(guardedAccesses, guardedAccess)
		case *ssa.Lookup:
			guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
			guardedAccesses = append(guardedAccesses, guardedAccess)
		case *ssa.Panic:
			guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
			guardedAccesses = append(guardedAccesses, guardedAccess)
		case *ssa.Range:
			guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
			guardedAccesses = append(guardedAccesses, guardedAccess)
		case *ssa.TypeAssert:
			guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
			guardedAccesses = append(guardedAccesses, guardedAccess)
		case *ssa.BinOp:
			guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
			guardedAccesses = append(guardedAccesses, guardedAccess)
			guardedAccess = domain.AddGuardedAccess(call.Pos(), call.Y, domain.GuardAccessRead, GoroutineState)
			guardedAccesses = append(guardedAccesses, guardedAccess)
		case *ssa.If:
			guardedAccess := domain.AddGuardedAccess(call.Pos(), call.Cond, domain.GuardAccessRead, GoroutineState)
			guardedAccesses = append(guardedAccesses, guardedAccess)
		case *ssa.MapUpdate:
			guardedAccess := domain.AddGuardedAccess(call.Pos(), call.Map, domain.GuardAccessWrite, GoroutineState)
			guardedAccesses = append(guardedAccesses, guardedAccess)
			guardedAccess = domain.AddGuardedAccess(call.Pos(), call.Value, domain.GuardAccessRead, GoroutineState)
			guardedAccesses = append(guardedAccesses, guardedAccess)
		case *ssa.Store:
			guardedAccess := domain.AddGuardedAccess(call.Pos(), call.Val, domain.GuardAccessRead, GoroutineState)
			guardedAccesses = append(guardedAccesses, guardedAccess)
			guardedAccess = domain.AddGuardedAccess(call.Pos(), call.Addr, domain.GuardAccessWrite, GoroutineState)
			guardedAccesses = append(guardedAccesses, guardedAccess)
		case *ssa.Call:
			callCommon := call.Common()
			guardedAccessesRet, guardedState := GetSummary(callCommon, GoroutineState.Copy())
			guardedAccesses = append(guardedAccesses, guardedAccessesRet...)
			GoroutineState.MergeStates(guardedState, false)
		case *ssa.Go:
			callCommon := call.Common()
			GoroutineState.Increment()
			newState := domain.NewGoroutineState()
			newState.Clock = GoroutineState.Clock
			newState.Lockset = GoroutineState.Lockset
			guardedAccessesRet, _ := GetSummary(callCommon, newState.Copy())
			guardedAccesses = append(guardedAccesses, guardedAccessesRet...)
		case *ssa.Defer:
			callCommon := call.Common()
			deferredFunctions = append(deferredFunctions, callCommon)
		}
	}
	return deferredFunctions, guardedAccesses, GoroutineState
}

func HandleFunction(GoroutineState *domain.GoroutineState, fn *ssa.Function) ([]*domain.GuardedAccess, *domain.GoroutineState) {
	guardedAccesses := make([]*domain.GuardedAccess, 0)
	var conditionalBlocks = map[string]struct{}{
		"if.then":     {},
		"if.else":     {},
		"select.body": {},
	}

	pkgName := fn.Pkg.Pkg.Path() // Used to guard against entering standard library packages
	packageToCheck := utils.GetTopPackageName() // The top package of the code. Any function under it is ok.
	if !strings.Contains(pkgName, packageToCheck) {
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
		res, GoroutineStateRet := GetSummary(deferredFunctions[i].Function, GoroutineState.Copy())
		guardedAccesses = append(guardedAccesses, res...)
		GoroutineState.MergeStates(GoroutineStateRet, deferredFunctions[i].IsConditional)
	}
	return guardedAccesses, GoroutineState
}

func GetSummary(callCommon *ssa.CallCommon, GoroutineState *domain.GoroutineState) ([]*domain.GuardedAccess, *domain.GoroutineState) {
	guardedAccesses := make([]*domain.GuardedAccess, 0)
	switch call := callCommon.Value.(type) {
	case *ssa.Builtin:
		return HandleBuiltin(GoroutineState, call, callCommon.Args)
	case *ssa.MakeClosure:
		fn := callCommon.Value.(*ssa.MakeClosure).Fn.(*ssa.Function)
		return HandleFunction(GoroutineState, fn)
	case *ssa.Function:
		if utils.IsCallTo(call, "(*sync.Mutex).Lock") {
			receiver := callCommon.Args[0]
			LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
			locks := map[string]*ssa.CallCommon{LockName: callCommon}
			GoroutineState.Lockset.UpdateLockSet(locks, nil)
			return guardedAccesses, GoroutineState
		}
		if utils.IsCallTo(call, "(*sync.Mutex).Unlock") {
			receiver := callCommon.Args[0]
			LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
			locks := map[string]*ssa.CallCommon{LockName: callCommon}
			GoroutineState.Lockset.UpdateLockSet(nil, locks)
			return guardedAccesses, GoroutineState
		}
		return HandleFunction(GoroutineState, call)
	}
	return guardedAccesses, GoroutineState
}
