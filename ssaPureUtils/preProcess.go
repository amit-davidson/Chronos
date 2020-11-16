package ssaPureUtils

import (
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
	locksCount        map[int]int
	visitedFuncs      *stacks.FunctionStackWithMap
}

func InitFunctionWithLocksPreprocess(entryFunc *ssa.Function) *FunctionWithLocksPreprocess {
	preProcess := &FunctionWithLocksPreprocess{
		FunctionWithLocks: make(map[*types.Signature]bool),
		locksCount:        make(map[int]int),
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
		GlobalPackageName, retError = GetTopLevelPackageName(pkg)
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

// IsFunctionContainingLocks calculates if a function contains locks depending on the call graph. It does it by
// iterating the entire call graph and calculates on each function the mutexes being held at each exit point (locksCount).
func IsFunctionContainingLocks(FunctionWithLocksPreprocess *FunctionWithLocksPreprocess, f *ssa.Function) bool {
	FunctionWithLocksPreprocess.visitedFuncs.Push(f)
	defer FunctionWithLocksPreprocess.visitedFuncs.Pop()

	if IsLock(f) || IsUnlock(f) {
		return false
	}

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
				if _, ok := FunctionWithLocksPreprocess.FunctionWithLocks[f.Signature]; ok {
					continue
				}
				if FunctionWithLocksPreprocess.visitedFuncs.Contains(f) {
					continue
				}
				if IsLock(f) {
					recv := callCommon.Args[len(callCommon.Args)-1]
					mutexPos := int(GetMutexPos(recv))
					FunctionWithLocksPreprocess.locksCount[mutexPos]++
				}
				if IsUnlock(f) {
					recv := callCommon.Args[len(callCommon.Args)-1]
					mutexPos := int(GetMutexPos(recv))
					if FunctionWithLocksPreprocess.locksCount[mutexPos] > 0 {
						FunctionWithLocksPreprocess.locksCount[mutexPos]--
						if FunctionWithLocksPreprocess.locksCount[mutexPos] == 0 {
							delete(FunctionWithLocksPreprocess.locksCount, mutexPos)
						}
					}
				}

				IsFunctionContainingLocks(FunctionWithLocksPreprocess, f)
			}
		}
	}
	var res bool
	if len(FunctionWithLocksPreprocess.locksCount) > 0 {
		res = true
	} else {
		res = false
	}
	FunctionWithLocksPreprocess.FunctionWithLocks[f.Signature] = res
	return res
}
