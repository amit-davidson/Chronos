package ssaUtils

import (
	"fmt"
	"go/types"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

var GlobalProgram *ssa.Program

var typesCache = make(map[*types.Interface][]types.Type, 0)

func LoadPackage(path string) (*ssa.Program, *ssa.Package, error) {
	conf1 := packages.Config{
		Mode: packages.LoadAllSyntax,
	}
	loadQuery := fmt.Sprintf("file=%s", path)
	pkgs, err := packages.Load(&conf1, loadQuery)
	if err != nil {
		return nil, nil, err
	}
	ssaProg, ssaPkgs := ssautil.AllPackages(pkgs, 0)
	ssaProg.Build()
	ssaPkg := ssaPkgs[0]
	return ssaProg, ssaPkg, nil
}
func SetGlobalProgram(prog *ssa.Program) () {
	GlobalProgram = prog
}
//func MapConcreteToInterface() {
//	for _, typ1 := range GlobalProgram.RuntimeTypes() {
//		for _, typ2 := range GlobalProgram.RuntimeTypes() {
//			typ1Interface, isType1Interface := typ1.Underlying().(*types.Interface)
//			typ2Interface, isType2Interface := typ2.Underlying().(*types.Interface)
//			name1 := typ1.String()
//			_ = name1
//			name2 := typ2.String()
//			_ = name2
//			if strings.Contains(name1, "Jerry") && strings.Contains(name2, "IceCreamMaker") {
//				fmt.Print("a")
//			}
//			if isType1Interface && types.Implements(typ2, typ1Interface) {
//				typesCache[typ1Interface] = append(typesCache[typ1Interface], typ2)
//			}
//			if isType2Interface && types.Implements(typ1, typ2Interface) {
//				typesCache[typ2Interface] = append(typesCache[typ2Interface], typ1)
//			}
//		}
//	}
//}

func GetMethodImplementations(recv types.Type, method *types.Func) []*ssa.Function {
	implementors := make([]types.Type, 0)
	methodImplementations := make([]*ssa.Function, 0)
	recvInterface := recv.(*types.Interface)
	for _, typ1 := range GlobalProgram.RuntimeTypes() {
		if types.Implements(typ1, recvInterface) {
			implementors = append(implementors, typ1)
		}
	}
	typesCache[recvInterface] = implementors
	for _, implementor := range implementors {
		setMethods := GlobalProgram.MethodSets.MethodSet(implementor)
		structMethod := setMethods.Lookup(method.Pkg(), method.Name())
		methodImplementations = append(methodImplementations, GlobalProgram.MethodValue(structMethod))
	}
	return methodImplementations
}
