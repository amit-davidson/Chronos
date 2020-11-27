package main

import (
	"flag"
	"fmt"
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/output"
	"github.com/amit-davidson/Chronos/pointerAnalysis"
	"github.com/amit-davidson/Chronos/ssaUtils"
	"github.com/amit-davidson/Chronos/utils"
	"golang.org/x/tools/go/ssa"
	"os"
)

func main() {
	defaultFile := flag.String("file", "", "The file containing the entry point of the program")
	defaultModulePath := flag.String("mod", "", "Path to the module where the search should be performed. It needs to be in the format:{VCS}/{organization}/{package}. Packages outside this path are excluded rom the search.")
	flag.Parse()
	if *defaultFile == "" {
		fmt.Printf("Please provide a file to load\n")
		os.Exit(1)
	}
	if *defaultModulePath == "" {
		fmt.Printf("Please provide a path to the module. It should be in the following format:{VCS}/{organization}/{package}.\n")
		os.Exit(1)
	}
	domain.GoroutineCounter = utils.NewCounter()
	domain.GuardedAccessCounter = utils.NewCounter()
	domain.PosIDCounter = utils.NewCounter()

	ssaProg, ssaPkg, err := ssaUtils.LoadPackage(*defaultFile, *defaultModulePath)
	if err != nil {
		fmt.Printf("Failed loading with the following error:%s\n", err)
		os.Exit(1)
	}
	entryFunc := ssaPkg.Func("main")
	err = ssaUtils.InitPreProcess(ssaProg, *defaultModulePath)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	entryCallCommon := ssa.CallCommon{Value: entryFunc}
	functionState := ssaUtils.HandleCallCommon(domain.NewEmptyContext(), &entryCallCommon, entryFunc.Pos())
	conflictingGAs, err := pointerAnalysis.Analysis(ssaPkg, functionState.GuardedAccesses)
	if err != nil {
		fmt.Printf("Error in analysis:%s\n", err)
		os.Exit(1)
	}
	err = output.GenerateError(conflictingGAs, ssaProg)
	if err != nil {
		fmt.Printf("Error in generating errors:%s\n", err)
		os.Exit(1)
	}
}
