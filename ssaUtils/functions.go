package ssaUtils

import (
	"StaticRaceDetector/domain"
	"StaticRaceDetector/utils"
	"golang.org/x/tools/go/ssa"
	"strings"
)

func HandleCallCommon(GoroutineState *domain.GoroutineState, callCommon *ssa.CallCommon) *domain.FunctionState {
	funcState := domain.GetEmptyFunctionState()
	if callCommon.IsInvoke() { // abstract methods (of interfaces) aren't handled
		return funcState
	}

	switch call := callCommon.Value.(type) {
	case *ssa.Builtin:
		HandleBuiltin(funcState, GoroutineState, call, callCommon.Args)
		return funcState
	case *ssa.MakeClosure:
		fn := callCommon.Value.(*ssa.MakeClosure).Fn.(*ssa.Function)
		funcStateRet := HandleFunction(GoroutineState, fn)
		return funcStateRet
	case *ssa.Function:
		if utils.IsCallTo(call, "(*sync.Mutex).Lock") {
			AddLock(funcState, callCommon, false)
			return funcState
		}
		if utils.IsCallTo(call, "(*sync.Mutex).Unlock") {
			AddLock(funcState, callCommon, true)
			return funcState
		}
		funcStateRet := HandleFunction(GoroutineState, call)
		return funcStateRet

	case ssa.Instruction:
		HandleInstruction(funcState, GoroutineState, call)
		return funcState
	}
	return funcState
}

func HandleBuiltin(functionState *domain.FunctionState, GoroutineState *domain.GoroutineState, callCommon *ssa.Builtin, args []ssa.Value) {
	switch name := callCommon.Name(); name {
	case "delete":
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[0], domain.GuardAccessWrite, functionState.Lockset, GoroutineState))
	case "cap", "len":
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[0], domain.GuardAccessRead, functionState.Lockset, GoroutineState))
	case "append":
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[1], domain.GuardAccessRead, functionState.Lockset, GoroutineState))
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[0], domain.GuardAccessWrite, functionState.Lockset, GoroutineState))
	case "copy":
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[0], domain.GuardAccessRead, functionState.Lockset, GoroutineState))
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(callCommon.Pos(), args[1], domain.GuardAccessWrite, functionState.Lockset, GoroutineState))
	}
}

func HandleInstruction(functionState *domain.FunctionState, GoroutineState *domain.GoroutineState, ins ssa.Instruction) {
	switch call := ins.(type) {
	case *ssa.UnOp:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, GoroutineState)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Field:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, GoroutineState)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.FieldAddr:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, GoroutineState)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Index:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, GoroutineState)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.IndexAddr:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, GoroutineState)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Lookup:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, GoroutineState)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Panic:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, GoroutineState)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Range:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, GoroutineState)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.TypeAssert:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, GoroutineState)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.BinOp:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, GoroutineState)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
		guardedAccess = domain.AddGuardedAccess(call.Pos(), call.Y, domain.GuardAccessRead, functionState.Lockset, GoroutineState)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.If:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.Cond, domain.GuardAccessRead, functionState.Lockset, GoroutineState)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.MapUpdate:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.Map, domain.GuardAccessWrite, functionState.Lockset, GoroutineState)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
		guardedAccess = domain.AddGuardedAccess(call.Pos(), call.Value, domain.GuardAccessRead, functionState.Lockset, GoroutineState)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Store:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.Val, domain.GuardAccessRead, functionState.Lockset, GoroutineState)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
		guardedAccess = domain.AddGuardedAccess(call.Pos(), call.Addr, domain.GuardAccessWrite, functionState.Lockset, GoroutineState)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	}
}

func GetBlockSummary(GoroutineState *domain.GoroutineState, block *ssa.BasicBlock) *domain.FunctionState {
	funcState := domain.GetEmptyFunctionState()
	for _, ins := range block.Instrs {
		switch call := ins.(type) {
		case *ssa.Call:
			callCommon := call.Common()
			funcStateRet := HandleCallCommon(GoroutineState, callCommon)
			funcState.MergeStates(funcStateRet)
		case *ssa.Go:
			callCommon := call.Common()
			newState := domain.NewGoroutineExecutionState(GoroutineState)
			funcStateRet := HandleCallCommon(newState.Copy(), callCommon)
			funcState.MergeStatesAfterGoroutine(funcStateRet)
		case *ssa.Defer:
			callCommon := call.Common()
			deferFunction :=  &domain.DeferFunction{Function: callCommon, BlockIndex:block.Index}
			funcState.DeferredFunctions = append(funcState.DeferredFunctions, deferFunction)
		default:
			HandleInstruction(funcState, GoroutineState, ins)
		}
	}
	return funcState
}

func HandleFunction(GoroutineState *domain.GoroutineState, fn *ssa.Function) *domain.FunctionState {
	funcState := domain.GetEmptyFunctionState()
	pkgName := fn.Pkg.Pkg.Path()                // Used to guard against entering standard library packages
	packageToCheck := utils.GetTopPackageName() // The top package of the code. Any function under it is ok.
	if !strings.Contains(pkgName, packageToCheck) {
		return funcState
	}

	//deferredFunctions := make([]*domain.FunctionWithBlock, 0)

	//deferredNodesToBlockLocksets := make(map[int]*domain.Lockset, 0)
	cfg := newCFG()
	funcState = cfg.getBlocksSummaries(GoroutineState, fn.Blocks[0])
	deferredMap := make(map[int][]*domain.DeferFunction, 0)
	for _, block := range funcState.DeferredFunctions {
		deferredMap[block.BlockIndex] = append(deferredMap[block.BlockIndex], block)
	}
	cfg.DeferredFunctions = deferredMap
	funcStateDefers := cfg.getDefersSummaries(GoroutineState, cfg.lastBlock)
	funcState.MergeStates(funcStateDefers)
	return funcState
}
