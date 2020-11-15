package ssaUtils

import (
	"github.com/amit-davidson/Chronos/ssaPureUtils"
	"github.com/amit-davidson/Chronos/utils/stacks"
	"go/types"
	"golang.org/x/tools/go/ssa"
)

var isFunctionContainingLocksMap = make(map[*types.Signature]bool)

var visitedFuncs = stacks.NewFunctionStackWithMap()

func PreProcess(entryFunc *ssa.Function) {
	IsFunctionContainingLocks(entryFunc)
}

func IsFunctionContainingLocks(f *ssa.Function) bool {
	visitedFuncs.Push(f)
	defer visitedFuncs.Pop()

	if IsLock(f) || IsUnlock(f) {
		return false
	}

	locksCount := make(map[int]int)
	isSubFunctionContainingLocks := false
	functionName := f.Name()
	_ = functionName
	for _, block := range f.Blocks {
		for _, instr := range block.Instrs {
			call, ok := instr.(ssa.CallInstruction)
			if !ok {
				continue
			}

			var funcs []*ssa.Function
			callCommon := call.Common()
			if callCommon.IsInvoke() {
				funcs = GetMethodImplementations(callCommon.Value.Type().Underlying(), callCommon.Method)
			} else {
				fnInstr, ok := callCommon.Value.(*ssa.Function)
				if !ok {
					anonFn, ok := callCommon.Value.(*ssa.MakeClosure)
					if !ok {
						continue
					}
					fnInstr = anonFn.Fn.(*ssa.Function)
				}
				funcs = []*ssa.Function{fnInstr}
			}

			for _, f := range funcs {
				if isContainingLocks, ok := isFunctionContainingLocksMap[f.Signature]; ok {
					if isContainingLocks == true {
						isSubFunctionContainingLocks = true
					}
					continue
				}
				if visitedFuncs.Contains(f) {
					continue
				}
				if IsLock(f) {
					recv := callCommon.Args[len(callCommon.Args)-1]
					mutexPos := int(ssaPureUtils.GetMutexPos(recv))
					locksCount[mutexPos]++
				}
				if IsUnlock(f) {
					recv := callCommon.Args[len(callCommon.Args)-1]
					mutexPos := int(ssaPureUtils.GetMutexPos(recv))
					if locksCount[mutexPos] > 0 {
						locksCount[mutexPos]--
					} else if locksCount[mutexPos] == 0 {
						delete(locksCount, mutexPos)
					}
				}

				res := IsFunctionContainingLocks(f)
				if res == true {
					isSubFunctionContainingLocks = true
				}
			}
		}
	}
	var res bool
	if len(locksCount) > 0 || isSubFunctionContainingLocks {
		res = true
	} else {
		res = false
	}
	isFunctionContainingLocksMap[f.Signature] = res
	return res
}
