package ssaUtils

import (
	"errors"
	"fmt"
	"github.com/amit-davidson/Chronos/domain"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"sort"
	"strings"
	"testing"
)

var typesCache = make(map[*types.Interface][]*ssa.Function)
var GlobalProgram *ssa.Program
var GlobalPackageName string
var PreProcess *PreProcessResults

var ErrNoPackages = errors.New("no packages in the path")


func Create(t *testing.T, path, fileName string) *ssa.Package {
	var conf loader.Config
	f, err := conf.ParseFile(fileName, nil)
	if err != nil {
		t.Fatal(err)
	}
	conf.CreateFromFiles(path, f)

	lprog, err := conf.Load()
	if err != nil {
		t.Fatal(err)
	}

	// We needn't call Build.
	foo := lprog.Package(path).Pkg
	return ssautil.CreateProgram(lprog, ssa.SanityCheckFunctions).Package(foo)
}

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

func LoadTests(path string) (*ssa.Program, *ssa.Package, error) {
	conf1 := packages.Config{
		Mode:  packages.LoadAllSyntax,
		Tests: true,
	}
	pkgs, err := packages.Load(&conf1, path)
	if err != nil {
		return nil, nil, err
	}
	if len(pkgs) == 0 {
		return nil, nil, fmt.Errorf("%s: %w", path, ErrNoPackages)
	}
	ssaProg, ssaPkgs := ssautil.AllPackages(pkgs, 0)
	ssaProg.Build()
	ssaPkg := ssaPkgs[1]
	return ssaProg, ssaPkg, nil
}

func GetTests(prog *ssa.Program, pkg *ssa.Package) []*ssa.Function {
	tests := make([]*ssa.Function, 0)
	for _, mem := range pkg.Members {
		if f, ok := mem.(*ssa.Function); ok &&
			ast.IsExported(f.Name()) &&
			strings.HasSuffix(prog.Fset.Position(f.Pos()).Filename, "_test.go") {

			switch {
			case isTest(f.Name(), "Test"):
				tests = append(tests, f)
			default:
				continue
			}
		}
	}
	return tests
}

func isTest(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if len(name) == len(prefix) { // "Test" is ok
		return true
	}
	return ast.IsExported(name[len(prefix):])
}

func GetStackTrace(prog *ssa.Program, ga *domain.GuardedAccess) string {
	stack := ""
	for _, pos := range ga.State.StackTrace.Iter() {
		calculatedPos := prog.Fset.Position(token.Pos(pos))
		stack += calculatedPos.String()
		stack += " ->\n"
	}
	return stack
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
