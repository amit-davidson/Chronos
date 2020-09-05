package main

func main() {
	//var conf loader.Config
	//file, err := conf.ParseFile("myprog.go", myprog)
	//if err != nil {
	//	fmt.Print(err) // parse error
	//	return
	//}
	//
	//// Create single-file main package and import its dependencies.
	//conf.CreateFromFiles("main", file)
	//
	//iprog, err := conf.Load()
	//if err != nil {
	//	fmt.Print(err) // type error in some package
	//	return
	//}
	//
	//// Create SSA-form program representation.
	//prog := ssautil.CreateProgram(iprog, 0)
	//mainPkg := prog.Package(iprog.Created[0].Pkg)
	//
	//// Build SSA code for bodies of all functions in the whole program.
	//prog.Build()
	//
	//funcInterface5 := mainPkg.Func("context3")
	//lsRet, guardedAccessRet := GetFunctionSummary(funcInterface5)
	//print(lsRet)
	//_, _ = lsRet, guardedAccessRet
}
