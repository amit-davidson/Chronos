package ssaUtils

import (
	"errors"
	"github.com/amit-davidson/Chronos/utils/stacks"
	"go/types"
	"golang.org/x/tools/go/ssa"
	"os"
	"path"
	"strings"
)

type FunctionWithLocksPreprocess struct {
	FunctionWithLocks map[*types.Signature]bool
	locks             map[int]struct{} // Is lock exists, and if it's in a conditional path
	visitedFuncs      *stacks.FunctionStackWithMap
}

func InitPreProcess(prog *ssa.Program, defaultModulePath string) error {
	GlobalProgram = prog
	p := strings.TrimSuffix(defaultModulePath, string(os.PathSeparator))
	splittedPath := strings.Split(p, string(os.PathSeparator))
	if len(splittedPath) < 3 {
		return errors.New("package should be provided in the following format:{VCS}/{organization}/{package}")
	}
	l := len(splittedPath)
	moduleName := path.Join(splittedPath[l-3], splittedPath[l-2], splittedPath[l-1])
	GlobalModuleName = moduleName
	return nil
}
