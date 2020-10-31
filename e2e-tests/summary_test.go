package e2e_tests

import (
	"encoding/json"
	"fmt"
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/e2e-tests/testutils"
	"github.com/amit-davidson/Chronos/pointerAnalysis"
	"github.com/amit-davidson/Chronos/ssaUtils"
	"github.com/amit-davidson/Chronos/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/ssa"
	"os"
	"testing"
)

const (
	updateAll = iota
	notUpdateAll
	testSelective
)

var shouldUpdateAll = testSelective

func TestE2E(t *testing.T) {
	var testCases = []struct {
		name         string
		testPath     string
		resPath      string
		shouldUpdate bool
	}{
		{name: "Lock", testPath: "locksAndUnlocks/Lock/prog1.go", resPath: "locksAndUnlocks/Lock/prog1_expected.json", shouldUpdate: false},
		{name: "LockAndUnlock", testPath: "locksAndUnlocks/LockAndUnlock/prog1.go", resPath: "locksAndUnlocks/LockAndUnlock/prog1_expected.json", shouldUpdate: false},
		{name: "LockAndUnlockIfBranch", testPath: "locksAndUnlocks/LockAndUnlockIfBranch/prog1.go", resPath: "locksAndUnlocks/LockAndUnlockIfBranch/prog1_expected.json", shouldUpdate: false},
		{name: "ElseIf", testPath: "general/ElseIf/prog1.go", resPath: "general/ElseIf/prog1_expected.json", shouldUpdate: false},
		{name: "LockInBothBranches", testPath: "locksAndUnlocks/LockInBothBranches/prog1.go", resPath: "locksAndUnlocks/LockInBothBranches/prog1_expected.json", shouldUpdate: false},
		{name: "DeferredLockAndUnlockIfBranch", testPath: "defer/DeferredLockAndUnlockIfBranch/prog1.go", resPath: "defer/DeferredLockAndUnlockIfBranch/prog1_expected.json", shouldUpdate: false},
		{name: "NestedDeferWithLockAndUnlock", testPath: "defer/NestedDeferWithLockAndUnlock/prog1.go", resPath: "defer/NestedDeferWithLockAndUnlock/prog1_expected.json", shouldUpdate: false},
		{name: "NestedDeferWithLockAndUnlockAndGoroutine", testPath: "defer/NestedDeferWithLockAndUnlockAndGoroutine/prog1.go", resPath: "defer/NestedDeferWithLockAndUnlockAndGoroutine/prog1_expected.json", shouldUpdate: false},
		{name: "LockAndUnlockIfMap", testPath: "locksAndUnlocks/LockAndUnlockIfMap/prog1.go", resPath: "locksAndUnlocks/LockAndUnlockIfMap/prog1_expected.json", shouldUpdate: false},
		{name: "MultipleLocksNoRace", testPath: "locksAndUnlocks/MultipleLocksNoRace/prog1.go", resPath: "locksAndUnlocks/MultipleLocksNoRace/prog1_expected.json", shouldUpdate: false},
		{name: "MultipleLocksRace", testPath: "locksAndUnlocks/MultipleLocksRace/prog1.go", resPath: "locksAndUnlocks/MultipleLocksRace/prog1_expected.json", shouldUpdate: false},
		{name: "NestedFunctions", testPath: "general/NestedFunctions/prog1.go", resPath: "general/NestedFunctions/prog1_expected.json", shouldUpdate: false},
		{name: "Simple", testPath: "general/Simple/prog1.go", resPath: "general/Simple/prog1_expected.json", shouldUpdate: false},
		{name: "NestedConditionWithLockInAllBranches", testPath: "locksAndUnlocks/NestedConditionWithLockInAllBranches/prog1.go", resPath: "locksAndUnlocks/NestedConditionWithLockInAllBranches/prog1_expected.json", shouldUpdate: false},
		{name: "DataRaceMap", testPath: "general/DataRaceMap/prog1.go", resPath: "general/DataRaceMap/prog1_expected.json", shouldUpdate: false},
		{name: "ForLoop", testPath: "unsupported/ForLoop/prog1.go", resPath: "unsupported/ForLoop/prog1_expected.json", shouldUpdate: false},
		{name: "DataRaceShadowedErr", testPath: "general/DataRaceShadowedErr/prog1.go", resPath: "general/DataRaceShadowedErr/prog1_expected.json", shouldUpdate: false},
		{name: "DataRaceInterfaceOverChannel", testPath: "pointerAnalysis/DataRaceInterfaceOverChannel/prog1.go", resPath: "pointerAnalysis/DataRaceInterfaceOverChannel/prog1_expected.json", shouldUpdate: false},
		{name: "DataRaceProperty", testPath: "general/DataRaceProperty/prog1.go", resPath: "general/DataRaceProperty/prog1_expected.json", shouldUpdate: false},
		{name: "DataRaceWithOnlyAlloc", testPath: "general/DataRaceWithOnlyAlloc/prog1.go", resPath: "general/DataRaceWithOnlyAlloc/prog1_expected.json", shouldUpdate: false},
		{name: "LockInsideGoroutine", testPath: "locksAndUnlocks/LockInsideGoroutine/prog1.go", resPath: "locksAndUnlocks/LockInsideGoroutine/prog1_expected.json", shouldUpdate: false},
		{name: "DataRaceWithSameFunction", testPath: "general/DataRaceWithSameFunction/prog1.go", resPath: "general/DataRaceWithSameFunction/prog1_expected.json", shouldUpdate: false},
		{name: "StructMethod", testPath: "general/StructMethod/prog1.go", resPath: "general/StructMethod/prog1_expected.json", shouldUpdate: false},
		{name: "DataRaceIceCreamMaker", testPath: "interfaces/DataRaceIceCreamMaker/prog1.go", resPath: "interfaces/DataRaceIceCreamMaker/prog1_expected.json", shouldUpdate: false},
		{name: "InterfaceWithLock", testPath: "interfaces/InterfaceWithLock/prog1.go", resPath: "interfaces/InterfaceWithLock/prog1_expected.json", shouldUpdate: false},
		{name: "NestedInterface", testPath: "interfaces/NestedInterface/prog1.go", resPath: "interfaces/NestedInterface/prog1_expected.json", shouldUpdate: false},
		{name: "TestNoRaceStackPushPop", testPath: "stdlib/TestNoRaceStackPushPop/prog1.go", resPath: "stdlib/TestNoRaceStackPushPop/prog1_expected.json", shouldUpdate: false},
		{name: "RaceNestedArrayCopy", testPath: "stdlib/RaceNestedArrayCopy/prog1.go", resPath: "stdlib/RaceNestedArrayCopy/prog1_expected.json", shouldUpdate: false},
		{name: "TestNoRaceAsFunc4", testPath: "stdlib/TestNoRaceAsFunc4/prog1.go", resPath: "stdlib/TestNoRaceAsFunc4/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceAsFunc3", testPath: "stdlib/TestRaceAsFunc3/prog1.go", resPath: "stdlib/TestRaceAsFunc3/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceAsFunc2", testPath: "stdlib/TestRaceAsFunc2/prog1.go", resPath: "stdlib/TestRaceAsFunc2/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceAsFunc1", testPath: "stdlib/TestRaceAsFunc1/prog1.go", resPath: "stdlib/TestRaceAsFunc1/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceCaseTypeIssue5890", testPath: "stdlib/TestRaceCaseTypeIssue5890/prog1.go", resPath: "stdlib/TestRaceCaseTypeIssue5890/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceCaseIssue6418", testPath: "stdlib/TestRaceCaseIssue6418/prog1.go", resPath: "stdlib/TestRaceCaseIssue6418/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceCaseFallthrough", testPath: "stdlib/TestRaceCaseFallthrough/prog1.go", resPath: "stdlib/TestRaceCaseFallthrough/prog1_expected.json", shouldUpdate: false},
		//{name: "TestNoRaceBlank", testPath: "stdlibNoSuccess/TestNoRaceBlank/prog1.go", resPath: "stdlibNoSuccess/TestNoRaceBlank/prog1_expected.json", shouldUpdate: false}, // blank space
		//{name: "TestRaceMethodThunk4", testPath: "stdlibNoSuccess/TestRaceMethodThunk4/prog1.go", resPath: "stdlibNoSuccess/TestRaceMethodThunk4/prog1_expected.json", shouldUpdate: false}, // Might be a bug in pointer analysis
		{name: "TestRaceMethodThunk3", testPath: "stdlib/TestRaceMethodThunk3/prog1.go", resPath: "stdlib/TestRaceMethodThunk3/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceMethodThunk2", testPath: "stdlib/TestRaceMethodThunk2/prog1.go", resPath: "stdlib/TestRaceMethodThunk2/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceMethodThunk", testPath: "stdlibNoSuccess/TestRaceMethodThunk/prog1.go", resPath: "stdlibNoSuccess/TestRaceMethodThunk/prog1_expected.json", shouldUpdate: false}, // blank space
		{name: "TestNoRaceMethodThunk", testPath: "stdlib/TestNoRaceMethodThunk/prog1.go", resPath: "stdlib/TestNoRaceMethodThunk/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceNestedStruct", testPath: "stdlib/TestRaceNestedStruct/prog1.go", resPath: "stdlib/TestRaceNestedStruct/prog1_expected.json", shouldUpdate: false},
		{name: "TestNoRaceEmptyStruct", testPath: "stdlibNoSuccess/TestNoRaceEmptyStruct/prog1.go", resPath: "stdlibNoSuccess/TestNoRaceEmptyStruct/prog1_expected.json", shouldUpdate: false},
		//{name: "TestRaceHeapParam", testPath: "stdlib/TestRaceHeapParam/prog1.go", resPath: "stdlib/TestRaceHeapParam/prog1_expected.json", shouldUpdate: false}, // No ssa param as value. Might be a bug in ssa.
		{name: "TestRaceStructInd", testPath: "stdlib/TestRaceStructInd/prog1.go", resPath: "stdlib/TestRaceStructInd/prog1_expected.json", shouldUpdate: false},
		//{name: "TestRaceAppendSliceStruct", testPath: "stdlibNoSuccess/TestRaceAppendSliceStruct/prog1.go", resPath: "stdlibNoSuccess/TestRaceAppendSliceStruct/prog1_expected.json", shouldUpdate: false}, // spread operator can't tell which item are affected
		{name: "TestRaceSliceStruct", testPath: "stdlibNoSuccess/TestRaceSliceStruct/prog1.go", resPath: "stdlibNoSuccess/TestRaceSliceStruct/prog1_expected.json", shouldUpdate: false}, // same
		{name: "TestRaceSliceString", testPath: "stdlib/TestRaceSliceString/prog1.go", resPath: "stdlib/TestRaceSliceString/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceSliceSlice2", testPath: "stdlib/TestRaceSliceSlice2/prog1.go", resPath: "stdlib/TestRaceSliceSlice2/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceSliceSlice", testPath: "stdlib/TestRaceSliceSlice/prog1.go", resPath: "stdlib/TestRaceSliceSlice/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceBlockAs", testPath: "stdlib/TestRaceBlockAs/prog1.go", resPath: "stdlib/TestRaceBlockAs/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceTypeAssert", testPath: "stdlib/TestRaceTypeAssert/prog1.go", resPath: "stdlib/TestRaceTypeAssert/prog1_expected.json", shouldUpdate: false},
		//{name: "TestRaceAddrExpr", testPath: "stdlib/TestRaceAddrExpr/prog1.go", resPath: "stdlib/TestRaceAddrExpr/prog1_expected.json", shouldUpdate: false}, // Due to the way ssa works, it's not possible to differ between p.x = 5 and p{x:5}. The first option might participate in a data race but the second never.
		{name: "TestNoRaceAddrExpr", testPath: "stdlib/TestNoRaceAddrExpr/prog1.go", resPath: "stdlib/TestNoRaceAddrExpr/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceDeferArg2", testPath: "stdlib/TestRaceDeferArg2/prog1.go", resPath: "stdlib/TestRaceDeferArg2/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceDeferArg", testPath: "stdlib/TestRaceDeferArg/prog1.go", resPath: "stdlib/TestRaceDeferArg/prog1_expected.json", shouldUpdate: false},
		{name: "TestRacePanicArg", testPath: "stdlib/TestRacePanicArg/prog1.go", resPath: "stdlib/TestRacePanicArg/prog1_expected.json", shouldUpdate: false},
		{name: "TestNoRaceMethodValue", testPath: "stdlib/TestNoRaceMethodValue/prog1.go", resPath: "stdlib/TestNoRaceMethodValue/prog1_expected.json", shouldUpdate: false},
		//{name: "TestRaceMethodValue3", testPath: "stdlib/TestRaceMethodValue3/prog1.go", resPath: "stdlib/TestRaceMethodValue3/prog1_expected.json", shouldUpdate: false}, // Might be a bug in pointer analysis
		{name: "TestRaceMethodValue2", testPath: "stdlib/TestRaceMethodValue2/prog1.go", resPath: "stdlib/TestRaceMethodValue2/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceMethodValue", testPath: "stdlib/TestRaceMethodValue/prog1.go", resPath: "stdlib/TestRaceMethodValue/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceMethodCall2", testPath: "stdlib/TestRaceMethodCall2/prog1.go", resPath: "stdlib/TestRaceMethodCall2/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceMethodCall", testPath: "stdlib/TestRaceMethodCall/prog1.go", resPath: "stdlib/TestRaceMethodCall/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceFuncCall", testPath: "stdlib/TestRaceFuncCall/prog1.go", resPath: "stdlib/TestRaceFuncCall/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceInterCall2", testPath: "stdlib/TestRaceInterCall2/prog1.go", resPath: "stdlib/TestRaceInterCall2/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceInterCall", testPath: "stdlib/TestRaceInterCall/prog1.go", resPath: "stdlib/TestRaceInterCall/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceMapInit2", testPath: "stdlib/TestRaceMapInit2/prog1.go", resPath: "stdlib/TestRaceMapInit2/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceMapInit", testPath: "stdlib/TestRaceMapInit/prog1.go", resPath: "stdlib/TestRaceMapInit/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceArrayInit", testPath: "stdlib/TestRaceArrayInit/prog1.go", resPath: "stdlib/TestRaceArrayInit/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceStructInit", testPath: "stdlib/TestRaceStructInit/prog1.go", resPath: "stdlib/TestRaceStructInit/prog1_expected.json", shouldUpdate: false},
		//{name: "TestNoRaceFuncUnlock", testPath: "stdlibNoSuccess/TestNoRaceFuncUnlock/prog1.go", resPath: "stdlibNoSuccess/TestNoRaceFuncUnlock/prog1_expected.json", shouldUpdate: false}, // No pointer analysis for locks
		{name: "TestRaceFuncItself", testPath: "stdlib/TestRaceFuncItself/prog1.go", resPath: "stdlib/TestRaceFuncItself/prog1_expected.json", shouldUpdate: false},
		//{name: "TestNoRaceShortCalc2", testPath: "stdlibNoSuccess/TestNoRaceShortCalc2/prog1.go", resPath: "stdlibNoSuccess/TestNoRaceShortCalc2/prog1_expected.json", shouldUpdate: false}, // Cant evaluate first part of condition to see the second will never execute
		//{name: "TestNoRaceShortCalc", testPath: "stdlibNoSuccess/TestNoRaceShortCalc/prog1.go", resPath: "stdlibNoSuccess/TestNoRaceShortCalc/prog1_expected.json", shouldUpdate: false}, // Cant evaluate first part of condition to see the second will never execute
		//{name: "TestNoRaceOr", testPath: "stdlibNoSuccess/TestNoRaceOr/prog1.go", resPath: "stdlibNoSuccess/TestNoRaceOr/prog1_expected.json", shouldUpdate: false}, // Cant evaluate first part of condition to see the second will never execute
		//{name: "TestRaceOr2", testPath: "stdlib/TestRaceOr2/prog1.go", resPath: "stdlib/TestRaceOr2/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceOr", testPath: "stdlib/TestRaceOr/prog1.go", resPath: "stdlib/TestRaceOr/prog1_expected.json", shouldUpdate: false},
		{name: "TestNoRaceAnd", testPath: "stdlibNoSuccess/TestNoRaceAnd/prog1.go", resPath: "stdlibNoSuccess/TestNoRaceAnd/prog1_expected.json", shouldUpdate: false}, // Cant evaluate first part of condition to see the second will never execute
		{name: "TestRaceAnd2", testPath: "stdlib/TestRaceAnd2/prog1.go", resPath: "stdlib/TestRaceAnd2/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceAnd", testPath: "stdlib/TestRaceAnd/prog1.go", resPath: "stdlib/TestRaceAnd/prog1_expected.json", shouldUpdate: false},
		//{name: "TestRaceEmptyInterface2", testPath: "stdlibNoSuccess/TestRaceEmptyInterface2/prog1.go", resPath: "stdlibNoSuccess/TestRaceEmptyInterface2/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceEmptyInterface1", testPath: "stdlib/TestRaceEmptyInterface/prog1.go", resPath: "stdlib/TestRaceEmptyInterface/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceRune", testPath: "stdlib/TestRaceRune/prog1.go", resPath: "stdlib/TestRaceRune/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceIndirection", testPath: "stdlibNoSuccess/TestRaceIndirection/prog1.go", resPath: "stdlibNoSuccess/TestRaceIndirection/prog1_expected.json", shouldUpdate: false}, // sync using channels
		{name: "TestRaceFuncArgsRW", testPath: "stdlib/TestRaceFuncArgsRW/prog1.go", resPath: "stdlib/TestRaceFuncArgsRW/prog1_expected.json", shouldUpdate: false},
		{name: "TestNoRaceFuncArgsRW", testPath: "stdlib/TestNoRaceFuncArgsRW/prog1.go", resPath: "stdlib/TestNoRaceFuncArgsRW/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceAppendCapRW", testPath: "stdlib/TestRaceAppendCapRW/prog1.go", resPath: "stdlib/TestRaceAppendCapRW/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceAppendLenRW", testPath: "stdlib/TestRaceAppendLenRW/prog1.go", resPath: "stdlib/TestRaceAppendLenRW/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceAppendRW", testPath: "stdlib/TestRaceAppendRW/prog1.go", resPath: "stdlib/TestRaceAppendRW/prog1_expected.json", shouldUpdate: false},
		//{name: "TestRacePanic", testPath: "stdlibNoSuccess/TestRacePanic/prog1.go", resPath: "stdlibNoSuccess/TestRacePanic/prog1_expected.json", shouldUpdate: false}, // cfg is weird because of the recover
		{name: "TestRaceFuncVariableWW", testPath: "stdlib/TestRaceFuncVariableWW/prog1.go", resPath: "stdlib/TestRaceFuncVariableWW/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceFuncVariableRW", testPath: "stdlib/TestRaceFuncVariableRW/prog1.go", resPath: "stdlib/TestRaceFuncVariableRW/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceUnsafePtrRW", testPath: "stdlib/TestRaceUnsafePtrRW/prog1.go", resPath: "stdlib/TestRaceUnsafePtrRW/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceComplex128WW", testPath: "stdlib/TestRaceComplex128WW/prog1.go", resPath: "stdlib/TestRaceComplex128WW/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceFloat64WW", testPath: "stdlib/TestRaceFloat64WW/prog1.go", resPath: "stdlib/TestRaceFloat64WW/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceStringPtrRW", testPath: "stdlib/TestRaceStringPtrRW/prog1.go", resPath: "stdlib/TestRaceStringPtrRW/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceStringRW", testPath: "stdlib/TestRaceStringRW/prog1.go", resPath: "stdlib/TestRaceStringRW/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceIntptrRW", testPath: "stdlib/TestRaceIntptrRW/prog1.go", resPath: "stdlib/TestRaceIntptrRW/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceError", testPath: "stdlib/TestRaceError/prog1.go", resPath: "stdlib/TestRaceError/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceIfaceConv", testPath: "stdlib/TestRaceIfaceConv/prog1.go", resPath: "stdlib/TestRaceIfaceConv/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceEfaceConv", testPath: "stdlib/TestRaceEfaceConv/prog1.go", resPath: "stdlib/TestRaceEfaceConv/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceIfaceCmpNil", testPath: "stdlib/TestRaceIfaceCmpNil/prog1.go", resPath: "stdlib/TestRaceIfaceCmpNil/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceIfaceCmp", testPath: "stdlib/TestRaceIfaceCmp/prog1.go", resPath: "stdlib/TestRaceIfaceCmp/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceIfaceWW", testPath: "stdlib/TestRaceIfaceWW/prog1.go", resPath: "stdlib/TestRaceIfaceWW/prog1_expected.json", shouldUpdate: false}, // Before write, a read is performed. So the creation confused with the read later.
		{name: "TestRaceEfaceWW", testPath: "stdlib/TestRaceEfaceWW/prog1.go", resPath: "stdlib/TestRaceEfaceWW/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceStructFieldRW3", testPath: "stdlib/TestRaceStructFieldRW3/prog1.go", resPath: "stdlib/TestRaceStructFieldRW3/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceStructFieldRW2", testPath: "stdlib/TestRaceStructFieldRW2/prog1.go", resPath: "stdlib/TestRaceStructFieldRW2/prog1_expected.json", shouldUpdate: false},
		{name: "TestNoRaceStructFieldRW2", testPath: "stdlib/TestNoRaceStructFieldRW2/prog1.go", resPath: "stdlib/TestNoRaceStructFieldRW2/prog1_expected.json", shouldUpdate: false},
		{name: "TestNoRaceStructFieldRW1", testPath: "stdlib/TestNoRaceStructFieldRW1/prog1.go", resPath: "stdlib/TestNoRaceStructFieldRW1/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceStructFieldRW1", testPath: "stdlib/TestRaceStructFieldRW1/prog1.go", resPath: "stdlib/TestRaceStructFieldRW1/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceStructRW", testPath: "stdlibNoSuccess/TestRaceStructRW/prog1.go", resPath: "stdlibNoSuccess/TestRaceStructRW/prog1_expected.json", shouldUpdate: false}, // The compiler optimizes the ssa in a way that instead of allocating on line 16, the fields are modified. It means that according to ssa, there shouldn't be any data race since the end result is the same. Maybe a bug in ssa?
		{name: "TestRaceArrayCopy", testPath: "stdlib/TestRaceArrayCopy/prog1.go", resPath: "stdlib/TestRaceArrayCopy/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceSprint", testPath: "stdlib/TestRaceSprint/prog1.go", resPath: "stdlib/TestRaceSprint/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceFuncArgument2", testPath: "stdlib/TestRaceFuncArgument2/prog1.go", resPath: "stdlib/TestRaceFuncArgument2/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceFuncArgument", testPath: "stdlib/TestRaceFuncArgument/prog1.go", resPath: "stdlib/TestRaceFuncArgument/prog1_expected.json", shouldUpdate: false},
		{name: "TestNoRaceEnoughRegisters", testPath: "stdlib/TestNoRaceEnoughRegisters/prog1.go", resPath: "stdlib/TestNoRaceEnoughRegisters/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceRotate", testPath: "stdlib/TestRaceRotate/prog1.go", resPath: "stdlib/TestRaceRotate/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceModConst", testPath: "stdlib/TestRaceModConst/prog1.go", resPath: "stdlib/TestRaceModConst/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceMod", testPath: "stdlib/TestRaceMod/prog1.go", resPath: "stdlib/TestRaceMod/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceDivConst", testPath: "stdlib/TestRaceDivConst/prog1.go", resPath: "stdlib/TestRaceDivConst/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceDiv", testPath: "stdlib/TestRaceDiv/prog1.go", resPath: "stdlib/TestRaceDiv/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceComplement", testPath: "stdlib/TestRaceComplement/prog1.go", resPath: "stdlib/TestRaceComplement/prog1_expected.json", shouldUpdate: false},
		{name: "TestNoRacePlus", testPath: "stdlib/TestNoRacePlus/prog1.go", resPath: "stdlib/TestNoRacePlus/prog1_expected.json", shouldUpdate: false},
		{name: "TestRacePlus2", testPath: "stdlib/TestRacePlus2/prog1.go", resPath: "stdlib/TestRacePlus2/prog1_expected.json", shouldUpdate: false},
		{name: "TestRacePlus", testPath: "stdlib/TestRacePlus/prog1.go", resPath: "stdlib/TestRacePlus/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceCaseTypeBody", testPath: "stdlib/TestRaceCaseTypeBody/prog1.go", resPath: "stdlib/TestRaceCaseTypeBody/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceCaseType", testPath: "stdlib/TestRaceCaseType/prog1.go", resPath: "stdlib/TestRaceCaseType/prog1_expected.json", shouldUpdate: false},
		//{name: "TestNoRaceCaseFallthrough", testPath: "stdlibNoSuccess/TestNoRaceCaseFallthrough/prog1.go", resPath: "stdlibNoSuccess/TestNoRaceCaseFallthrough/prog1_expected.json", shouldUpdate: false}, // No way to determine flow as the detector is flow insensitive
		{name: "TestRaceCaseBody", testPath: "stdlib/TestRaceCaseBody/prog1.go", resPath: "stdlib/TestRaceCaseBody/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceCaseCondition2", testPath: "stdlib/TestRaceCaseCondition2/prog1.go", resPath: "stdlib/TestRaceCaseCondition2/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceCaseCondition", testPath: "stdlib/TestRaceCaseCondition/prog1.go", resPath: "stdlib/TestRaceCaseCondition/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceInt32RWClosures", testPath: "stdlib/TestRaceInt32RWClosures/prog1.go", resPath: "stdlib/TestRaceInt32RWClosures/prog1_expected.json", shouldUpdate: false},
		{name: "TestNoRaceIntRWClosures", testPath: "stdlib/TestNoRaceIntRWClosures/prog1.go", resPath: "stdlib/TestNoRaceIntRWClosures/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceIntRWClosures", testPath: "stdlib/TestRaceIntRWClosures/prog1.go", resPath: "stdlib/TestRaceIntRWClosures/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceIntRWGlobalFuncs", testPath: "stdlib/TestRaceIntRWGlobalFuncs/prog1.go", resPath: "stdlib/TestRaceIntRWGlobalFuncs/prog1_expected.json", shouldUpdate: false},
		{name: "TestNoRaceComp", testPath: "stdlib/TestNoRaceComp/prog1.go", resPath: "stdlib/TestNoRaceComp/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceComp2", testPath: "stdlib/TestRaceComp2/prog1.go", resPath: "stdlib/TestRaceComp2/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceSelect1", testPath: "stdlibNoSuccess/TestRaceSelect1/prog1.go", resPath: "stdlibNoSuccess/TestRaceSelect1/prog1_expected.json", shouldUpdate: false},
		{name: "TestRaceUnaddressableMapLen", testPath: "stdlib/TestRaceUnaddressableMapLen/prog1.go", resPath: "stdlib/TestRaceUnaddressableMapLen/prog1_expected.json", shouldUpdate: false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ssaProg, ssaPkg, err := ssaUtils.LoadPackage(tc.testPath)
			require.NoError(t, err)

			err = ssaUtils.SetGlobals(ssaProg, ssaPkg, "")
			require.NoError(t, err)


			entryFunc := ssaPkg.Func("main")

			domain.GoroutineCounter.Reset()
			domain.GuardedAccessCounter.Reset()
			entryCallCommon := ssa.CallCommon{Value: entryFunc}
			functionState := ssaUtils.HandleCallCommon(domain.NewEmptyContext(), &entryCallCommon, entryFunc.Pos())
			testresult := testutils.NewTestResult(functionState.Lockset, functionState.GuardedAccesses)
			dump, err := json.MarshalIndent(testresult, "", "\t")
			require.NoError(t, err)
			if shouldUpdateAll == updateAll || shouldUpdateAll == testSelective && tc.shouldUpdate {
				utils.UpdateFile(t, tc.resPath, dump)
			}
			testutils.CompareResult(t, tc.resPath, functionState.Lockset, functionState.GuardedAccesses)
			err = pointerAnalysis.Analysis(ssaPkg, ssaProg, functionState.GuardedAccesses)
			if err != nil {
				fmt.Printf("Error in analysis:%s\n", err)
				os.Exit(1)
			}
		})
	}
}
