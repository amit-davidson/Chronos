package e2e_tests

import (
	"fmt"
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/output"
	"github.com/amit-davidson/Chronos/pointerAnalysis"
	"github.com/amit-davidson/Chronos/ssaUtils"
	"github.com/amit-davidson/Chronos/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/ssa"
	"os"
	"testing"
)

func TestStdlib(t *testing.T) {
	var testCases = []struct {
		name     string
		testPath string
	}{
		{name: "TestNoRaceStackPushPop", testPath: "testdata/stdlib/TestNoRaceStackPushPop/prog1.go"},
		{name: "RaceNestedArrayCopy", testPath: "testdata/stdlib/RaceNestedArrayCopy/prog1.go"},
		{name: "TestNoRaceAsFunc4", testPath: "testdata/stdlib/TestNoRaceAsFunc4/prog1.go"},
		{name: "TestRaceAsFunc3", testPath: "testdata/stdlib/TestRaceAsFunc3/prog1.go"},
		{name: "TestRaceAsFunc2", testPath: "testdata/stdlib/TestRaceAsFunc2/prog1.go"},
		{name: "TestRaceAsFunc1", testPath: "testdata/stdlib/TestRaceAsFunc1/prog1.go"},
		{name: "TestRaceCaseTypeIssue5890", testPath: "testdata/stdlib/TestRaceCaseTypeIssue5890/prog1.go"},
		{name: "TestRaceCaseIssue6418", testPath: "testdata/stdlib/TestRaceCaseIssue6418/prog1.go"},
		{name: "TestRaceCaseFallthrough", testPath: "testdata/stdlib/TestRaceCaseFallthrough/prog1.go"},
		//{name: "TestNoRaceBlank", testPath: "testdata/stdlibNoSuccess/TestNoRaceBlank/prog1.go"},  // blank space
		//{name: "TestRaceMethodThunk4", testPath: "testdata/stdlibNoSuccess/TestRaceMethodThunk4/prog1.go"},  // Might be a bug in pointer analysis
		{name: "TestRaceMethodThunk3", testPath: "testdata/stdlib/TestRaceMethodThunk3/prog1.go"},
		{name: "TestRaceMethodThunk2", testPath: "testdata/stdlib/TestRaceMethodThunk2/prog1.go"},
		{name: "TestRaceMethodThunk", testPath: "testdata/stdlibNoSuccess/TestRaceMethodThunk/prog1.go"}, // blank space
		{name: "TestNoRaceMethodThunk", testPath: "testdata/stdlib/TestNoRaceMethodThunk/prog1.go"},
		{name: "TestRaceNestedStruct", testPath: "testdata/stdlib/TestRaceNestedStruct/prog1.go"},
		{name: "TestNoRaceEmptyStruct", testPath: "testdata/stdlibNoSuccess/TestNoRaceEmptyStruct/prog1.go"},
		//{name: "TestRaceHeapParam", testPath: "testdata/stdlib/TestRaceHeapParam/prog1.go"},  // No ssa param as value. Might be a bug in ssa.
		{name: "TestRaceStructInd", testPath: "testdata/stdlib/TestRaceStructInd/prog1.go"},
		//{name: "TestRaceAppendSliceStruct", testPath: "testdata/stdlibNoSuccess/TestRaceAppendSliceStruct/prog1.go"},  // spread operator can't tell which item are affected
		{name: "TestRaceSliceStruct", testPath: "testdata/stdlibNoSuccess/TestRaceSliceStruct/prog1.go"}, // same
		{name: "TestRaceSliceString", testPath: "testdata/stdlib/TestRaceSliceString/prog1.go"},
		{name: "TestRaceSliceSlice2", testPath: "testdata/stdlib/TestRaceSliceSlice2/prog1.go"},
		{name: "TestRaceSliceSlice", testPath: "testdata/stdlib/TestRaceSliceSlice/prog1.go"},
		{name: "TestRaceBlockAs", testPath: "testdata/stdlib/TestRaceBlockAs/prog1.go"},
		{name: "TestRaceTypeAssert", testPath: "testdata/stdlib/TestRaceTypeAssert/prog1.go"},
		//{name: "TestRaceAddrExpr", testPath: "testdata/stdlib/TestRaceAddrExpr/prog1.go"},  // Due to the way ssa works, it's not possible to differ between p.x = 5 and p{x:5}. The first option might participate in a data race but the second never.
		{name: "TestNoRaceAddrExpr", testPath: "testdata/stdlib/TestNoRaceAddrExpr/prog1.go"},
		{name: "TestRaceDeferArg2", testPath: "testdata/stdlib/TestRaceDeferArg2/prog1.go"},
		{name: "TestRaceDeferArg", testPath: "testdata/stdlib/TestRaceDeferArg/prog1.go"},
		{name: "TestRacePanicArg", testPath: "testdata/stdlib/TestRacePanicArg/prog1.go"},
		{name: "TestNoRaceMethodValue", testPath: "testdata/stdlib/TestNoRaceMethodValue/prog1.go"},
		//{name: "TestRaceMethodValue3", testPath: "testdata/stdlib/TestRaceMethodValue3/prog1.go"},  // Might be a bug in pointer analysis
		{name: "TestRaceMethodValue2", testPath: "testdata/stdlib/TestRaceMethodValue2/prog1.go"},
		{name: "TestRaceMethodValue", testPath: "testdata/stdlib/TestRaceMethodValue/prog1.go"},
		{name: "TestRaceMethodCall2", testPath: "testdata/stdlib/TestRaceMethodCall2/prog1.go"},
		{name: "TestRaceMethodCall", testPath: "testdata/stdlib/TestRaceMethodCall/prog1.go"},
		{name: "TestRaceFuncCall", testPath: "testdata/stdlib/TestRaceFuncCall/prog1.go"},
		{name: "TestRaceInterCall2", testPath: "testdata/stdlib/TestRaceInterCall2/prog1.go"},
		{name: "TestRaceInterCall", testPath: "testdata/stdlib/TestRaceInterCall/prog1.go"},
		{name: "TestRaceMapInit2", testPath: "testdata/stdlib/TestRaceMapInit2/prog1.go"},
		{name: "TestRaceMapInit", testPath: "testdata/stdlib/TestRaceMapInit/prog1.go"},
		{name: "TestRaceArrayInit", testPath: "testdata/stdlib/TestRaceArrayInit/prog1.go"},
		{name: "TestRaceStructInit", testPath: "testdata/stdlib/TestRaceStructInit/prog1.go"},
		//{name: "TestNoRaceFuncUnlock", testPath: "testdata/stdlibNoSuccess/TestNoRaceFuncUnlock/prog1.go"},  // No pointer analysis for locks
		{name: "TestRaceFuncItself", testPath: "testdata/stdlib/TestRaceFuncItself/prog1.go"},
		//{name: "TestNoRaceShortCalc2", testPath: "testdata/stdlibNoSuccess/TestNoRaceShortCalc2/prog1.go"},  // Cant evaluate first part of condition to see the second will never execute
		//{name: "TestNoRaceShortCalc", testPath: "testdata/stdlibNoSuccess/TestNoRaceShortCalc/prog1.go"},  // Cant evaluate first part of condition to see the second will never execute
		//{name: "TestNoRaceOr", testPath: "testdata/stdlibNoSuccess/TestNoRaceOr/prog1.go"},  // Cant evaluate first part of condition to see the second will never execute
		//{name: "TestRaceOr2", testPath: "testdata/stdlib/TestRaceOr2/prog1.go"},
		{name: "TestRaceOr", testPath: "testdata/stdlib/TestRaceOr/prog1.go"},
		{name: "TestNoRaceAnd", testPath: "testdata/stdlibNoSuccess/TestNoRaceAnd/prog1.go"}, // Cant evaluate first part of condition to see the second will never execute
		{name: "TestRaceAnd2", testPath: "testdata/stdlib/TestRaceAnd2/prog1.go"},
		{name: "TestRaceAnd", testPath: "testdata/stdlib/TestRaceAnd/prog1.go"},
		//{name: "TestRaceEmptyInterface2", testPath: "testdata/stdlibNoSuccess/TestRaceEmptyInterface2/prog1.go"},
		{name: "TestRaceEmptyInterface1", testPath: "testdata/stdlib/TestRaceEmptyInterface/prog1.go"},
		{name: "TestRaceRune", testPath: "testdata/stdlib/TestRaceRune/prog1.go"},
		{name: "TestRaceIndirection", testPath: "testdata/stdlibNoSuccess/TestRaceIndirection/prog1.go"}, // sync using channels
		{name: "TestRaceFuncArgsRW", testPath: "testdata/stdlib/TestRaceFuncArgsRW/prog1.go"},
		{name: "TestNoRaceFuncArgsRW", testPath: "testdata/stdlib/TestNoRaceFuncArgsRW/prog1.go"},
		{name: "TestRaceAppendCapRW", testPath: "testdata/stdlib/TestRaceAppendCapRW/prog1.go"},
		{name: "TestRaceAppendLenRW", testPath: "testdata/stdlib/TestRaceAppendLenRW/prog1.go"},
		{name: "TestRaceAppendRW", testPath: "testdata/stdlib/TestRaceAppendRW/prog1.go"},
		//{name: "TestRacePanic", testPath: "testdata/stdlibNoSuccess/TestRacePanic/prog1.go"},  // cfg is weird because of the recover
		{name: "TestRaceFuncVariableWW", testPath: "testdata/stdlib/TestRaceFuncVariableWW/prog1.go"},
		{name: "TestRaceFuncVariableRW", testPath: "testdata/stdlib/TestRaceFuncVariableRW/prog1.go"},
		{name: "TestRaceUnsafePtrRW", testPath: "testdata/stdlib/TestRaceUnsafePtrRW/prog1.go"},
		{name: "TestRaceComplex128WW", testPath: "testdata/stdlib/TestRaceComplex128WW/prog1.go"},
		{name: "TestRaceFloat64WW", testPath: "testdata/stdlib/TestRaceFloat64WW/prog1.go"},
		{name: "TestRaceStringPtrRW", testPath: "testdata/stdlib/TestRaceStringPtrRW/prog1.go"},
		{name: "TestRaceStringRW", testPath: "testdata/stdlib/TestRaceStringRW/prog1.go"},
		{name: "TestRaceIntptrRW", testPath: "testdata/stdlib/TestRaceIntptrRW/prog1.go"},
		{name: "TestRaceError", testPath: "testdata/stdlib/TestRaceError/prog1.go"},
		{name: "TestRaceIfaceConv", testPath: "testdata/stdlib/TestRaceIfaceConv/prog1.go"},
		{name: "TestRaceEfaceConv", testPath: "testdata/stdlib/TestRaceEfaceConv/prog1.go"},
		{name: "TestRaceIfaceCmpNil", testPath: "testdata/stdlib/TestRaceIfaceCmpNil/prog1.go"},
		{name: "TestRaceIfaceCmp", testPath: "testdata/stdlib/TestRaceIfaceCmp/prog1.go"},
		{name: "TestRaceIfaceWW", testPath: "testdata/stdlib/TestRaceIfaceWW/prog1.go"}, // Before write, a read is performed. So the creation confused with the read later.
		{name: "TestRaceEfaceWW", testPath: "testdata/stdlib/TestRaceEfaceWW/prog1.go"},
		{name: "TestRaceStructFieldRW3", testPath: "testdata/stdlib/TestRaceStructFieldRW3/prog1.go"},
		{name: "TestRaceStructFieldRW2", testPath: "testdata/stdlib/TestRaceStructFieldRW2/prog1.go"},
		{name: "TestNoRaceStructFieldRW2", testPath: "testdata/stdlib/TestNoRaceStructFieldRW2/prog1.go"},
		{name: "TestNoRaceStructFieldRW1", testPath: "testdata/stdlib/TestNoRaceStructFieldRW1/prog1.go"},
		{name: "TestRaceStructFieldRW1", testPath: "testdata/stdlib/TestRaceStructFieldRW1/prog1.go"},
		{name: "TestRaceStructRW", testPath: "testdata/stdlibNoSuccess/TestRaceStructRW/prog1.go"}, // The compiler optimizes the ssa in a way that instead of allocating on line 16, the fields are modified. It means that according to ssa, there shouldn't be any data race since the end result is the same. Maybe a bug in ssa?
		{name: "TestRaceArrayCopy", testPath: "testdata/stdlib/TestRaceArrayCopy/prog1.go"},
		{name: "TestRaceSprint", testPath: "testdata/stdlib/TestRaceSprint/prog1.go"},
		{name: "TestRaceFuncArgument2", testPath: "testdata/stdlib/TestRaceFuncArgument2/prog1.go"},
		{name: "TestRaceFuncArgument", testPath: "testdata/stdlib/TestRaceFuncArgument/prog1.go"},
		{name: "TestNoRaceEnoughRegisters", testPath: "testdata/stdlib/TestNoRaceEnoughRegisters/prog1.go"},
		{name: "TestRaceRotate", testPath: "testdata/stdlib/TestRaceRotate/prog1.go"},
		{name: "TestRaceModConst", testPath: "testdata/stdlib/TestRaceModConst/prog1.go"},
		{name: "TestRaceMod", testPath: "testdata/stdlib/TestRaceMod/prog1.go"},
		{name: "TestRaceDivConst", testPath: "testdata/stdlib/TestRaceDivConst/prog1.go"},
		{name: "TestRaceDiv", testPath: "testdata/stdlib/TestRaceDiv/prog1.go"},
		{name: "TestRaceComplement", testPath: "testdata/stdlib/TestRaceComplement/prog1.go"},
		{name: "TestNoRacePlus", testPath: "testdata/stdlib/TestNoRacePlus/prog1.go"},
		{name: "TestRacePlus2", testPath: "testdata/stdlib/TestRacePlus2/prog1.go"},
		{name: "TestRacePlus", testPath: "testdata/stdlib/TestRacePlus/prog1.go"},
		{name: "TestRaceCaseTypeBody", testPath: "testdata/stdlib/TestRaceCaseTypeBody/prog1.go"},
		{name: "TestRaceCaseType", testPath: "testdata/stdlib/TestRaceCaseType/prog1.go"},
		//{name: "TestNoRaceCaseFallthrough", testPath: "testdata/stdlibNoSuccess/TestNoRaceCaseFallthrough/prog1.go"},  // No way to determine flow as the detector is flow insensitive
		{name: "TestRaceCaseBody", testPath: "testdata/stdlib/TestRaceCaseBody/prog1.go"},
		{name: "TestRaceCaseCondition2", testPath: "testdata/stdlib/TestRaceCaseCondition2/prog1.go"},
		{name: "TestRaceCaseCondition", testPath: "testdata/stdlib/TestRaceCaseCondition/prog1.go"},
		{name: "TestRaceInt32RWClosures", testPath: "testdata/stdlib/TestRaceInt32RWClosures/prog1.go"},
		{name: "TestNoRaceIntRWClosures", testPath: "testdata/stdlib/TestNoRaceIntRWClosures/prog1.go"},
		{name: "TestRaceIntRWClosures", testPath: "testdata/stdlib/TestRaceIntRWClosures/prog1.go"},
		{name: "TestRaceIntRWGlobalFuncs", testPath: "testdata/stdlib/TestRaceIntRWGlobalFuncs/prog1.go"},
		{name: "TestNoRaceComp", testPath: "testdata/stdlib/TestNoRaceComp/prog1.go"},
		{name: "TestRaceComp2", testPath: "testdata/stdlib/TestRaceComp2/prog1.go"},
		{name: "TestRaceSelect1", testPath: "testdata/stdlib/TestRaceSelect1/prog1.go"},
		{name: "TestRaceSelect2", testPath: "testdata/stdlib/TestRaceSelect2/prog1.go"},
		{name: "TestRaceSelect3", testPath: "testdata/stdlib/TestRaceSelect3/prog1.go"},
		{name: "TestRaceSelect4", testPath: "testdata/stdlib/TestRaceSelect4/prog1.go"},
		{name: "TestRaceSelect5", testPath: "testdata/stdlib/TestRaceSelect5/prog1.go"},
		//{name: "TestNoRaceSelect1", testPath: "testdata/stdlibNoSuccess/TestNoRaceSelect1/prog1.go"},  // All of the no race use syncing with channels
		//{name: "TestNoRaceSelect2", testPath: "testdata/stdlibNoSuccess/TestNoRaceSelect2/prog1.go"},
		//{name: "TestNoRaceSelect3", testPath: "testdata/stdlibNoSuccess/TestNoRaceSelect3/prog1.go"},
		//{name: "TestNoRaceSelect4", testPath: "testdata/stdlibNoSuccess/TestNoRaceSelect4/prog1.go"},
		//{name: "TestNoRaceSelect5", testPath: "testdata/stdlibNoSuccess/TestNoRaceSelect5/prog1.go"},
		{name: "TestRaceUnaddressableMapLen", testPath: "testdata/stdlib/TestRaceUnaddressableMapLen/prog1.go"},
		{name: "TestNoRaceCase", testPath: "testdata/stdlib/TestNoRaceCase/prog1.go"},
		{name: "TestNoRaceRangeIssue5446", testPath: "testdata/stdlib/TestNoRaceRangeIssue5446/prog1.go"},
		//{name: "TestRaceRange", testPath: "testdata/stdlibNoSuccess/TestRaceRange/prog1.go"},  // doesn't work. There's a race between v on line 11 from iter i+1 and 16/18 form iter i and it's not handled properly.
		{name: "TestRaceForInit", testPath: "testdata/stdlib/TestRaceForInit/prog1.go"},
		//{name: "TestNoRaceForInit", testPath: "testdata/stdlibNoSuccess/TestNoRaceForInit/prog1.go"},  // flow analysis required
		{name: "TestRaceForTest", testPath: "testdata/stdlib/TestRaceForTest/prog1.go"},
		{name: "TestRaceForIncr", testPath: "testdata/stdlib/TestRaceForIncr/prog1.go"},
		{name: "TestNoRaceForIncr", testPath: "testdata/stdlibNoSuccess/TestNoRaceForIncr/prog1.go"}, // flow analysis required
		{name: "TestNoRaceHeapReallocation", testPath: "testdata/stdlib/TestNoRaceHeapReallocation/prog1.go"},
		//{name: "TestRaceIssue5567", testPath: "testdata/stdlib/TestRaceIssue5567/prog1.go"},  // There's write inside f.Read on b which is an external package
		//{name: "TestRaceIssue5654", testPath: "testdata/stdlib/TestRaceIssue5654/prog1.go"},  // PackageProblem
		{name: "TestNoRaceTinyAlloc", testPath: "testdata/stdlib/TestNoRaceTinyAlloc/prog1.go"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ssaProg, ssaPkg, err := ssaUtils.LoadPackage(tc.testPath)
			require.NoError(t, err)

			domain.GoroutineCounter = utils.NewCounter()
			domain.GuardedAccessCounter = utils.NewCounter()
			domain.PosIDCounter = utils.NewCounter()

			entryFunc := ssaPkg.Func("main")
			ssaUtils.InitPreProcess(ssaProg, "github.com/amit-davidson/Chronos")

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
		})
	}
}
