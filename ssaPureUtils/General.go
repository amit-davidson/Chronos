package ssaPureUtils

import (
	"errors"
	"go/token"
	"golang.org/x/tools/go/ssa"
	"os"
	"strings"
)

var GlobalProgram *ssa.Program
var GlobalPackageName string



func GetMutexPos(value ssa.Value) token.Pos {
	val, ok := GetField(value)
	if !ok {
		return value.Pos()
	}
	obj := GetUnderlyingObjectFromField(val)
	return obj.Pos()

}

func SetGlobals(prog *ssa.Program, pkg *ssa.Package, defaultPkgPath string) error {
	GlobalProgram = prog
	if defaultPkgPath != "" {
		GlobalPackageName = strings.TrimSuffix(defaultPkgPath, string(os.PathSeparator))

		return nil
	}

	var retError error
	GlobalPackageName, retError = GetTopLevelPackageName(pkg)
	if retError != nil {
		return retError
	}
	return nil
}

func GetTopLevelPackageName(pkg *ssa.Package) (string, error) {
	pkgName := pkg.Pkg.Path()
	r := strings.SplitAfterN(pkgName, string(os.PathSeparator), 4)
	if len(r) < 3 {
		return "", errors.New("package should be provided in the following format:{VCS}/{organization}/{package}")
	}
	topLevelPackage := r[0] + r[1] + r[2]
	return topLevelPackage, nil
}