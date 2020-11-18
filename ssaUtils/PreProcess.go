package ssaUtils

import (
	"github.com/amit-davidson/Chronos/ssaPureUtils"
	"github.com/amit-davidson/Chronos/utils/stacks"
	"go/types"
	"golang.org/x/tools/go/ssa"
	"os"
	"strings"
)

type PreProcessResults struct {
	*FunctionWithLocksPreprocess
}

type FunctionWithLocksPreprocess struct {
	FunctionWithLocks map[*types.Signature]bool
	locks             map[int]struct{} // Is lock exists, and if it's in a conditional path
	visitedFuncs      *stacks.FunctionStackWithMap
}

func InitFunctionWithLocksPreprocess(entryFunc *ssa.Function) *FunctionWithLocksPreprocess {
	preProcess := &FunctionWithLocksPreprocess{
		FunctionWithLocks: make(map[*types.Signature]bool),
		locks:             make(map[int]struct{}),
		visitedFuncs:      stacks.NewFunctionStackWithMap(),
	}
	IsFunctionContainingLocks(preProcess, entryFunc)
	return preProcess
}
func InitPreProcess(prog *ssa.Program, pkg *ssa.Package, defaultPkgPath string, entryFunc *ssa.Function) error {
	GlobalProgram = prog
	if defaultPkgPath != "" {
		GlobalPackageName = strings.TrimSuffix(defaultPkgPath, string(os.PathSeparator))
	} else {
		var retError error
		GlobalPackageName, retError = ssaPureUtils.GetTopLevelPackageName(pkg)
		if retError != nil {
			return retError
		}
	}

	PreProcessResults := &PreProcessResults{}
	locksPreProcess := InitFunctionWithLocksPreprocess(entryFunc)
	PreProcessResults.FunctionWithLocksPreprocess = locksPreProcess
	PreProcess = PreProcessResults
	return nil
}

// IsFunctionContainingLocks calculates if a function contains locks when it finishes it's execution depending on the
// call graph. It does it by iterating the entire call graph and calculates on each block if a lock is succeeded with an
// unlock. If it does not, the function is marked as containing locks.
func IsFunctionContainingLocks(FunctionWithLocksPreprocess *FunctionWithLocksPreprocess, f *ssa.Function) bool {
	FunctionWithLocksPreprocess.visitedFuncs.Push(f)
	defer FunctionWithLocksPreprocess.visitedFuncs.Pop()

	if ssaPureUtils.IsLock(f) || ssaPureUtils.IsUnlock(f) {
		return false
	}

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
				if _, ok := FunctionWithLocksPreprocess.FunctionWithLocks[f.Signature]; ok {
					continue
				}
				if FunctionWithLocksPreprocess.visitedFuncs.Contains(f) {
					continue
				}
				if ssaPureUtils.IsLock(f) {
					recv := callCommon.Args[len(callCommon.Args)-1]
					mutexPos := int(ssaPureUtils.GetMutexPos(recv))
					FunctionWithLocksPreprocess.locks[mutexPos] = struct{}{}
				}
				if ssaPureUtils.IsUnlock(f) {
					recv := callCommon.Args[len(callCommon.Args)-1]
					mutexPos := int(ssaPureUtils.GetMutexPos(recv))
					_, isLockExist := FunctionWithLocksPreprocess.locks[mutexPos]
					if !isLockExist {
						continue
					}
					delete(FunctionWithLocksPreprocess.locks, mutexPos)
				}

				IsFunctionContainingLocks(FunctionWithLocksPreprocess, f)
			}
		}
	}
	var res bool
	if len(FunctionWithLocksPreprocess.locks) > 0 {
		res = true
	} else {
		res = false
	}
	FunctionWithLocksPreprocess.FunctionWithLocks[f.Signature] = res
	return res
}
