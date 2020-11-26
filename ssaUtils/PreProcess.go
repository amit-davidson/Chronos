package ssaUtils

import (
	"github.com/amit-davidson/Chronos/utils/stacks"
	"go/types"
	"golang.org/x/tools/go/ssa"
)

type FunctionWithLocksPreprocess struct {
	FunctionWithLocks map[*types.Signature]bool
	locks             map[int]struct{} // Is lock exists, and if it's in a conditional path
	visitedFuncs      *stacks.FunctionStackWithMap
}

func InitPreProcess(prog *ssa.Program, defaultModulePath string) {
	GlobalProgram = prog
	GlobalModuleName = defaultModulePath
}