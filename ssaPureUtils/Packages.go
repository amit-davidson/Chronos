package ssaPureUtils

import (
	"errors"
	"fmt"
	"go/types"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"sort"
)


var typesCache = make(map[*types.Interface][]*ssa.Function)

var ErrNoPackages = errors.New("no packages in the path")

func LoadPackage(path string) (*ssa.Program, *ssa.Package, error) {
	conf1 := packages.Config{
		Mode: packages.LoadAllSyntax,
	}
	loadQuery := fmt.Sprintf("file=%s", path)
	pkgs, err := packages.Load(&conf1, loadQuery)
	if err != nil {
		return nil, nil, err
	}
	if len(pkgs) == 0 {
		return nil, nil, fmt.Errorf("%s: %w", path, ErrNoPackages)
	}
	ssaProg, ssaPkgs := ssautil.AllPackages(pkgs, 0)
	ssaProg.Build()
	ssaPkg := ssaPkgs[0]
	return ssaProg, ssaPkg, nil
}

func GetMethodImplementations(recv types.Type, method *types.Func) []*ssa.Function {
	methodImplementations := make([]*ssa.Function, 0)
	recvInterface := recv.(*types.Interface)

	if methodImplementations, ok := typesCache[recvInterface]; ok {
		return methodImplementations
	}

	implementors := make([]types.Type, 0)
	for _, typ := range GlobalProgram.RuntimeTypes() {
		if types.Implements(typ, recvInterface) {
			implementors = append(implementors, typ)
		}
	}
	for _, implementor := range implementors {
		setMethods := GlobalProgram.MethodSets.MethodSet(implementor)
		method := setMethods.Lookup(method.Pkg(), method.Name())
		methodImpl := GlobalProgram.MethodValue(method)
		if methodImpl.Synthetic == "" {
			methodImplementations = append(methodImplementations, methodImpl)
		}
	}

	// Sort by pos to enter previous implementations first. This make the search deterministic and easier for debugging
	sortedImplementations := sortMethodImplementations(methodImplementations)
	typesCache[recvInterface] = sortedImplementations
	return sortedImplementations
}

func sortMethodImplementations(methodImplementations []*ssa.Function) []*ssa.Function {
	posSlice := make([]int, 0)
	sortedImplementations := make([]*ssa.Function, 0)
	implMap := make(map[int]*ssa.Function)
	for _, methodImplementation := range methodImplementations {
		pos := methodImplementation.Pos()
		implMap[int(pos)] = methodImplementation
		posSlice = append(posSlice, int(pos))
	}
	sort.Ints(posSlice)
	for _, pos := range posSlice {
		sortedImplementations = append(sortedImplementations, implMap[pos])
	}
	return sortedImplementations
}

