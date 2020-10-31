package ssaUtils

import (
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/utils"
	"github.com/amit-davidson/Chronos/utils/stacks"
	"go/token"
	"golang.org/x/tools/go/ssa"
	"strings"
)

func HandleCallCommon(Context *domain.Context, callCommon *ssa.CallCommon, pos token.Pos) *domain.FunctionState {
	Context.StackTrace.Push(int(pos))
	defer Context.StackTrace.Pop()

	funcState := domain.GetEmptyFunctionState()
	if callCommon.IsInvoke() {
		impls := GetMethodImplementations(callCommon.Value.Type().Underlying(), callCommon.Method)
		if len(impls) > 0 {
			funcState = HandleFunction(Context, impls[0])
			for _, impl := range impls[1:] {
				funcstateRet := HandleFunction(Context, impl)
				funcState.MergeBranchState(funcstateRet)
			}
		}
		return funcState
	}

	switch call := callCommon.Value.(type) {
	case *ssa.Builtin:
		HandleBuiltin(funcState, Context, callCommon)
		return funcState
	case *ssa.MakeClosure:
		fn := callCommon.Value.(*ssa.MakeClosure).Fn.(*ssa.Function)
		funcStateRet := HandleFunction(Context, fn)
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
		funcStateRet := HandleFunction(Context, call)
		return funcStateRet

	case ssa.Instruction:
		HandleInstruction(funcState, Context, call)
		return funcState
	}
	return funcState
}

func HandleBuiltin(functionState *domain.FunctionState, Context *domain.Context, call *ssa.CallCommon) {
	callCommon := call.Value.(*ssa.Builtin)
	args := call.Args
	switch name := callCommon.Name(); name {
	case "delete":
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(call.Pos(), args[0], domain.GuardAccessWrite, functionState.Lockset, Context))
	case "cap", "len":
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(call.Pos(), args[0], domain.GuardAccessRead, functionState.Lockset, Context))
	case "append":
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(call.Pos(), args[1], domain.GuardAccessRead, functionState.Lockset, Context))
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(call.Pos(), args[0], domain.GuardAccessWrite, functionState.Lockset, Context))
	case "copy":
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(call.Pos(), args[0], domain.GuardAccessRead, functionState.Lockset, Context))
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(call.Pos(), args[1], domain.GuardAccessWrite, functionState.Lockset, Context))
	}
}

func HandleInstruction(functionState *domain.FunctionState, Context *domain.Context, ins ssa.Instruction) {
	switch call := ins.(type) {
	case *ssa.UnOp:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, Context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Field:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call, domain.GuardAccessRead, functionState.Lockset, Context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.FieldAddr:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call, domain.GuardAccessRead, functionState.Lockset, Context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Index:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, Context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.IndexAddr:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, Context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Lookup:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, Context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Panic:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, Context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Range:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, Context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.TypeAssert:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, Context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.BinOp:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, Context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
		guardedAccess = domain.AddGuardedAccess(call.Pos(), call.Y, domain.GuardAccessRead, functionState.Lockset, Context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.If:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.Cond, domain.GuardAccessRead, functionState.Lockset, Context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.MapUpdate:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.Map, domain.GuardAccessWrite, functionState.Lockset, Context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
		guardedAccess = domain.AddGuardedAccess(call.Pos(), call.Value, domain.GuardAccessRead, functionState.Lockset, Context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Store:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.Val, domain.GuardAccessRead, functionState.Lockset, Context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
		guardedAccess = domain.AddGuardedAccess(call.Pos(), call.Addr, domain.GuardAccessWrite, functionState.Lockset, Context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Return:
		for _, retValue := range call.Results {
			guardedAccess := domain.AddGuardedAccess(call.Pos(), retValue, domain.GuardAccessRead, functionState.Lockset, Context)
			functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
		}
	}
}

func GetBlockSummary(Context *domain.Context, block *ssa.BasicBlock) *domain.FunctionState {
	funcState := domain.GetEmptyFunctionState()
	for _, ins := range block.Instrs {
		switch call := ins.(type) {
		case *ssa.Call:
			callCommon := call.Common()
			funcStateRet := HandleCallCommon(Context, callCommon, callCommon.Pos())
			funcState.MergeStates(funcStateRet, true)
		case *ssa.Go:
			callCommon := call.Common()
			newState := domain.NewGoroutineExecutionState(Context)
			funcStateRet := HandleCallCommon(newState.Copy(), callCommon, callCommon.Pos())
			funcState.MergeStates(funcStateRet, false)
		case *ssa.Defer:
			callCommon := call.Common()
			funcState.DeferredFunctions.Push(callCommon)
		default:
			HandleInstruction(funcState, Context, ins)
		}
	}
	return funcState
}

func (cfg *CFG) runDefers(Context *domain.Context, defers *stacks.FunctionStack) *domain.FunctionState {
	calculatedState := domain.GetEmptyFunctionState()
	for {
		deferFunction := defers.Pop()
		if deferFunction == nil {
			break
		}
		retState := HandleCallCommon(Context, deferFunction, deferFunction.Pos())
		retState.UpdateGuardedAccessesWithLockset(calculatedState.Lockset)
		calculatedState.MergeStates(retState, true)
	}
	return calculatedState

}

func HandleFunction(Context *domain.Context, fn *ssa.Function) *domain.FunctionState {
	funcState := domain.GetEmptyFunctionState()
	if fn.Pkg == nil {
		return funcState
	}
	pkgName := fn.Pkg.Pkg.Path() // Used to guard against entering standard library packages
	if !strings.Contains(pkgName, GlobalPackageName) {
		return funcState
	}

	// regular
	cfg := newCFG()
	cfg.calculateState(Context, fn.Blocks[0])
	return cfg.calculatedState
}
