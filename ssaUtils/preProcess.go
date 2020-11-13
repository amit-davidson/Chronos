package ssaUtils

import (
	"github.com/amit-davidson/Chronos/utils/stacks"
	"go/types"
	"golang.org/x/tools/go/ssa"
)

var isFunctionContainingLocksMap = make(map[*types.Signature]bool)

var visitedFuncs = stacks.NewFunctionStackWithMap()

func PreProcess(entryFunc *ssa.Function) {
	IsFunctionContainingLocks(entryFunc)
	print(5)
}

func IsFunctionContainingLocks(f *ssa.Function) {
	visitedFuncs.Push(f)
	defer visitedFuncs.Pop()

	functionName := f.Name()
	_ = functionName
	for _, block := range f.Blocks {
		for _, instr := range block.Instrs {
			call, ok := instr.(ssa.CallInstruction)
			if !ok {
				continue
			}

			callCommon := call.Common()
			if callCommon.IsInvoke() {
				impls := GetMethodImplementations(callCommon.Value.Type().Underlying(), callCommon.Method)
				for _, impl := range impls {
					if visitedFuncs.Contains(impl) {
						continue
					}
					IsFunctionContainingLocks(impl)
				}
			} else {
				fnInstr, ok := callCommon.Value.(*ssa.Function)
				if !ok {
					anonFn, ok := callCommon.Value.(*ssa.MakeClosure)
					if !ok {
						continue
					}
					fnInstr = anonFn.Fn.(*ssa.Function)
				}
				if isContainingLocks, ok := isFunctionContainingLocksMap[fnInstr.Signature]; ok {
					if isContainingLocks == true {
						isFunctionContainingLocksMap[f.Signature] = true
						return
					}
					continue
				}
				if visitedFuncs.Contains(fnInstr) {
					continue
				}
				IsFunctionContainingLocks(fnInstr)

				if IsLock(fnInstr) {
					isFunctionContainingLocksMap[f.Signature] = true
				}
			}
		}
		isFunctionContainingLocksMap[f.Signature] = false
	}
}
