package ssaUtils

import (
	"StaticRaceDetector/domain"
	"StaticRaceDetector/utils"
	"golang.org/x/tools/go/ssa"
	"strconv"
	"strings"
)

func GetSummary(guardedAccesses *[]*domain.GuardedAccess, GoroutineState *domain.GoroutineState, callCommon *ssa.CallCommon) *domain.GoroutineState {
	if callCommon.IsInvoke() { // abstract methods (of interfaces) aren't handled
		return GoroutineState
	}

	switch call := callCommon.Value.(type) {
	case *ssa.Builtin:
		return HandleBuiltin(guardedAccesses, GoroutineState, call, callCommon.Args)
	case *ssa.MakeClosure:
		fn := callCommon.Value.(*ssa.MakeClosure).Fn.(*ssa.Function)
		return HandleFunction(guardedAccesses, GoroutineState, fn)
	case *ssa.Function:
		if utils.IsCallTo(call, "(*sync.Mutex).Lock") {
			receiver := callCommon.Args[0]
			LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
			locks := map[string]*ssa.CallCommon{LockName: callCommon}
			GoroutineState.Lockset.UpdateLockSet(locks, nil)
			return GoroutineState
		}
		if utils.IsCallTo(call, "(*sync.Mutex).Unlock") {
			receiver := callCommon.Args[0]
			LockName := receiver.Name() + strconv.Itoa(int(receiver.Pos()))
			locks := map[string]*ssa.CallCommon{LockName: callCommon}
			GoroutineState.Lockset.UpdateLockSet(nil, locks)
			return GoroutineState
		}
		return HandleFunction(guardedAccesses, GoroutineState, call)
	case ssa.Instruction:
		HandleInstruction(guardedAccesses, GoroutineState, call)
		return GoroutineState
	}
	return GoroutineState
}

func HandleBuiltin(guardedAccesses *[]*domain.GuardedAccess, GoroutineState *domain.GoroutineState, callCommon *ssa.Builtin, args []ssa.Value) *domain.GoroutineState {
	switch name := callCommon.Name(); name {
	case "delete":
		*guardedAccesses = append(*guardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[0], domain.GuardAccessWrite, GoroutineState))
		return GoroutineState
	case "cap", "len":
		*guardedAccesses = append(*guardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[0], domain.GuardAccessRead, GoroutineState))
		return GoroutineState
	case "append":
		*guardedAccesses = append(*guardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[1], domain.GuardAccessRead, GoroutineState))
		*guardedAccesses = append(*guardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[0], domain.GuardAccessWrite, GoroutineState))
		return GoroutineState
	case "copy":
		*guardedAccesses = append(*guardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[0], domain.GuardAccessRead, GoroutineState))
		*guardedAccesses = append(*guardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[1], domain.GuardAccessWrite, GoroutineState))
		return GoroutineState
	}
	return GoroutineState
}

func HandleInstruction(guardedAccesses *[]*domain.GuardedAccess, GoroutineState *domain.GoroutineState, ins ssa.Instruction) {
	switch call := ins.(type) {
	case *ssa.UnOp:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		*guardedAccesses = append(*guardedAccesses, guardedAccess)
	case *ssa.Field:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		*guardedAccesses = append(*guardedAccesses, guardedAccess)
	case *ssa.FieldAddr:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		*guardedAccesses = append(*guardedAccesses, guardedAccess)
	case *ssa.Index:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		*guardedAccesses = append(*guardedAccesses, guardedAccess)
	case *ssa.IndexAddr:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		*guardedAccesses = append(*guardedAccesses, guardedAccess)
	case *ssa.Lookup:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		*guardedAccesses = append(*guardedAccesses, guardedAccess)
	case *ssa.Panic:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		*guardedAccesses = append(*guardedAccesses, guardedAccess)
	case *ssa.Range:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		*guardedAccesses = append(*guardedAccesses, guardedAccess)
	case *ssa.TypeAssert:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		*guardedAccesses = append(*guardedAccesses, guardedAccess)
	case *ssa.BinOp:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, GoroutineState)
		*guardedAccesses = append(*guardedAccesses, guardedAccess)
		guardedAccess = domain.AddGuardedAccess(call.Pos(), call.Y, domain.GuardAccessRead, GoroutineState)
		*guardedAccesses = append(*guardedAccesses, guardedAccess)
	case *ssa.If:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.Cond, domain.GuardAccessRead, GoroutineState)
		*guardedAccesses = append(*guardedAccesses, guardedAccess)
	case *ssa.MapUpdate:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.Map, domain.GuardAccessWrite, GoroutineState)
		*guardedAccesses = append(*guardedAccesses, guardedAccess)
		guardedAccess = domain.AddGuardedAccess(call.Pos(), call.Value, domain.GuardAccessRead, GoroutineState)
		*guardedAccesses = append(*guardedAccesses, guardedAccess)
	case *ssa.Store:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.Val, domain.GuardAccessRead, GoroutineState)
		*guardedAccesses = append(*guardedAccesses, guardedAccess)
		guardedAccess = domain.AddGuardedAccess(call.Pos(), call.Addr, domain.GuardAccessWrite, GoroutineState)
		*guardedAccesses = append(*guardedAccesses, guardedAccess)
	}
}

func GetBlockSummary(guardedAccesses *[]*domain.GuardedAccess, GoroutineState *domain.GoroutineState, block *ssa.BasicBlock) ([]*ssa.CallCommon, *domain.GoroutineState) {
	deferredFunctions := make([]*ssa.CallCommon, 0)
	for _, ins := range block.Instrs {
		switch call := ins.(type) {
		case *ssa.Call:
			callCommon := call.Common()
			guardedState := GetSummary(guardedAccesses, GoroutineState.Copy(), callCommon)
			GoroutineState.MergeStates(guardedState, false)
		case *ssa.Go:
			callCommon := call.Common()
			newState := domain.NewGoroutineExecutionState(GoroutineState)
			_ = GetSummary(guardedAccesses, newState.Copy(), callCommon)
		case *ssa.Defer:
			callCommon := call.Common()
			deferredFunctions = append(deferredFunctions, callCommon)
		default:
			HandleInstruction(guardedAccesses, GoroutineState, ins)
		}
	}
	return deferredFunctions, GoroutineState
}

func HandleFunction(guardedAccesses *[]*domain.GuardedAccess, GoroutineState *domain.GoroutineState, fn *ssa.Function) *domain.GoroutineState {
	var conditionalBlocks = map[string]struct{}{
		"if.then":     {},
		"if.else":     {},
		"select.body": {},
	}

	pkgName := fn.Pkg.Pkg.Path()                // Used to guard against entering standard library packages
	packageToCheck := utils.GetTopPackageName() // The top package of the code. Any function under it is ok.
	if !strings.Contains(pkgName, packageToCheck) {
		return GoroutineState
	}

	deferredFunctions := make([]*domain.ConditionalFunction, 0)
	for _, block := range fn.Blocks {
		deferredFunctionsRet, goroutineState := GetBlockSummary(guardedAccesses, GoroutineState.Copy(), block)

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
		GoroutineStateRet := GetSummary(guardedAccesses, GoroutineState.Copy(), deferredFunctions[i].Function)
		GoroutineState.MergeStates(GoroutineStateRet, deferredFunctions[i].IsConditional)
	}
	return GoroutineState
}
