package ssaUtils

import (
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/utils/stacks"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ssa"
	"strings"
)

var functionsCache = make(map[*types.Signature]*domain.FunctionState)

func HandleCallCommon(context *domain.Context, callCommon *ssa.CallCommon, pos token.Pos) *domain.BlockState {
	funcState := domain.GetEmptyBlockState()

	// if we already visited this path, it means we're (probably) in a recursion so we return to avoid infinite loop
	if context.StackTrace.Contains(int(pos)) {
		return funcState
	}

	context.StackTrace.Push(int(pos))
	defer context.StackTrace.Pop()

	if callCommon.IsInvoke() {
		impls := GetMethodImplementations(callCommon.Value.Type().Underlying(), callCommon.Method)
		if len(impls) > 0 {
			funcState = HandleFunction(context, impls[0])
			for _, impl := range impls[1:] {
				funcstateRet := HandleFunction(context, impl)
				funcState.MergeSiblingBlock(funcstateRet)
			}
		}
		return funcState
	}

	switch call := callCommon.Value.(type) {
	case *ssa.Builtin:
		HandleBuiltin(funcState, context, callCommon)
		return funcState
	case *ssa.MakeClosure:
		fn := callCommon.Value.(*ssa.MakeClosure).Fn.(*ssa.Function)
		funcStateRet := HandleFunction(context, fn)
		return funcStateRet
	case *ssa.Function:
		if IsLock(call) {
			AddLock(funcState, callCommon, false)
			return funcState
		}
		if IsUnlock(call) {
			AddLock(funcState, callCommon, true)
			return funcState
		}

		var blockStateRet *domain.BlockState
		sig := callCommon.Signature()
		if cachedFunctionState, ok := functionsCache[sig]; ok {
			copiedState := cachedFunctionState.Copy() // Copy to avoid override cached item
			copiedState.AddContextToFunction(context)
			blockStateRet = domain.CreateBlockState(copiedState.GuardedAccesses, copiedState.Lockset, stacks.NewCallCommonStack())
		} else {
			blockStateRet = HandleFunction(context, call)
			fs := domain.CreateFunctionState(blockStateRet.GuardedAccesses, blockStateRet.Lockset)
			fs.RemoveContextFromFunction(context)
			functionsCache[sig] = fs
		}
		return blockStateRet

	case ssa.Instruction:
		HandleInstruction(funcState, context, call)
		return funcState
	}
	return funcState
}

func HandleBuiltin(functionState *domain.BlockState, context *domain.Context, call *ssa.CallCommon) {
	callCommon := call.Value.(*ssa.Builtin)
	args := call.Args
	switch name := callCommon.Name(); name {
	case "delete":
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(call.Pos(), args[0], domain.GuardAccessWrite, functionState.Lockset, context))
	case "cap", "len":
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(call.Pos(), args[0], domain.GuardAccessRead, functionState.Lockset, context))
	case "append":
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(call.Pos(), args[1], domain.GuardAccessRead, functionState.Lockset, context))
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(call.Pos(), args[0], domain.GuardAccessWrite, functionState.Lockset, context))
	case "copy":
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(call.Pos(), args[0], domain.GuardAccessRead, functionState.Lockset, context))
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, domain.AddGuardedAccess(call.Pos(), args[1], domain.GuardAccessWrite, functionState.Lockset, context))
	}
}

func HandleInstruction(functionState *domain.BlockState, context *domain.Context, ins ssa.Instruction) {
	switch call := ins.(type) {
	case *ssa.UnOp:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Field:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call, domain.GuardAccessRead, functionState.Lockset, context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.FieldAddr:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call, domain.GuardAccessRead, functionState.Lockset, context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Index:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.IndexAddr:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Lookup:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Panic:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Range:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.TypeAssert:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.BinOp:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.X, domain.GuardAccessRead, functionState.Lockset, context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
		guardedAccess = domain.AddGuardedAccess(call.Pos(), call.Y, domain.GuardAccessRead, functionState.Lockset, context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.If:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.Cond, domain.GuardAccessRead, functionState.Lockset, context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.MapUpdate:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.Map, domain.GuardAccessWrite, functionState.Lockset, context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
		guardedAccess = domain.AddGuardedAccess(call.Pos(), call.Value, domain.GuardAccessRead, functionState.Lockset, context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Store:
		guardedAccess := domain.AddGuardedAccess(call.Pos(), call.Val, domain.GuardAccessRead, functionState.Lockset, context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
		guardedAccess = domain.AddGuardedAccess(call.Pos(), call.Addr, domain.GuardAccessWrite, functionState.Lockset, context)
		functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
	case *ssa.Return:
		for _, retValue := range call.Results {
			guardedAccess := domain.AddGuardedAccess(call.Pos(), retValue, domain.GuardAccessRead, functionState.Lockset, context)
			functionState.GuardedAccesses = append(functionState.GuardedAccesses, guardedAccess)
		}
	}
}

func GetBlockSummary(context *domain.Context, block *ssa.BasicBlock) *domain.BlockState {
	funcState := domain.GetEmptyBlockState()
	for _, ins := range block.Instrs {
		switch call := ins.(type) {
		case *ssa.Call:
			callCommon := call.Common()
			funcStateRet := HandleCallCommon(context, callCommon, callCommon.Pos())
			funcState.AddResult(funcStateRet, true)
		case *ssa.Go:
			callCommon := call.Common()
			newState := domain.NewGoroutineExecutionState(context)
			funcStateRet := HandleCallCommon(newState.Copy(), callCommon, callCommon.Pos())
			funcState.AddResult(funcStateRet, false)
		case *ssa.Defer:
			callCommon := call.Common()
			funcState.DeferredFunctions.Push(callCommon)
		default:
			HandleInstruction(funcState, context, ins)
		}
	}
	return funcState
}

func (cfg *CFG) runDefers(context *domain.Context, defers *stacks.CallCommonStack) *domain.BlockState {
	calculatedState := domain.GetEmptyBlockState()
	for {
		deferFunction := defers.Pop()
		if deferFunction == nil {
			break
		}
		retState := HandleCallCommon(context, deferFunction, deferFunction.Pos())
		calculatedState.MergeChildBlock(retState)
	}
	return calculatedState

}

func (cfg *CFG) calculateFunctionStatePathInsensitive(context *domain.Context, blocks []*ssa.BasicBlock) {
	for _, block := range blocks {
		cfg.calculateBlockStateIfNeeded(context, block)
		state := cfg.ComputedBlocks[block.Index].Copy()
		if cfg.calculatedState == nil {
			cfg.calculatedState = state
		} else {
			cfg.calculatedState.MergeChildBlock(state)
		}

		deferBlock := cfg.ComputedDeferBlocks[block.Index]
		if deferBlock != nil {
			cfg.calculatedState.MergeChildBlock(deferBlock)
		}
	}
}

func HandleFunction(context *domain.Context, fn *ssa.Function) *domain.BlockState {
	funcState := domain.GetEmptyBlockState()
	if fn.Pkg == nil {
		return funcState
	}
	pkgName := fn.Pkg.Pkg.Path() // Used to guard against entering standard library packages
	if !strings.Contains(pkgName, GlobalPackageName) {
		return funcState
	}

	// regular
	cfg := newCFG()
	if fn.Blocks == nil { // External function
		return funcState
	}
	isContainingLocks, ok := isFunctionContainingLocksMap[fn.Signature]
	if !ok {
		panic("Function is being iterated but wasn't found when iterating on program functions in preprocess")
	}
	if isContainingLocks {
		cfg.calculateFunctionStatePathSensitive(context, fn.Blocks[0])
	} else {
		cfg.calculateFunctionStatePathInsensitive(context, fn.Blocks)
	}

	return cfg.calculatedState
}
