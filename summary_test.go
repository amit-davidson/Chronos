package main

import (
	"StaticRaceDetector/domain"
	"StaticRaceDetector/ssaUtils"
	"StaticRaceDetector/testutils"
	"StaticRaceDetector/utils"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/ssa"
	"testing"
)

const (
	updateAll = iota
	notUpdateAll
	testSelective
)

var shouldUpdateAll = testSelective

func TestGetFunctionSummary(t *testing.T) {
	var testCases = []struct {
		name         string
		testPath     string
		resPath      string
		shouldUpdate bool
	}{
		//{name: "Lock", testPath: "e2e-tests/locksAndUnlocks/Lock/prog1.go", resPath: "e2e-tests/locksAndUnlocks/Lock/prog1_expected.json", shouldUpdate: false},
		//{name: "LockAndUnlock", testPath: "e2e-tests/locksAndUnlocks/LockAndUnlock/prog1.go", resPath: "e2e-tests/locksAndUnlocks/LockAndUnlock/prog1_expected.json", shouldUpdate: false},
		//{name: "LockAndUnlockIfBranch", testPath: "e2e-tests/locksAndUnlocks/LockAndUnlockIfBranch/prog1.go", resPath: "e2e-tests/locksAndUnlocks/LockAndUnlockIfBranch/prog1_expected.json", shouldUpdate: false},
		//{name: "ElseIf", testPath: "e2e-tests/general/ElseIf/prog1.go", resPath: "e2e-tests/general/ElseIf/prog1_expected.json", shouldUpdate: false},
		//{name: "LockInBothBranches", testPath: "e2e-tests/locksAndUnlocks/LockInBothBranches/prog1.go", resPath: "e2e-tests/locksAndUnlocks/LockInBothBranches/prog1_expected.json", shouldUpdate: false},
		//{name: "DeferredLockAndUnlockIfBranch", testPath: "e2e-tests/defer/DeferredLockAndUnlockIfBranch/prog1.go", resPath: "e2e-tests/defer/DeferredLockAndUnlockIfBranch/prog1_expected.json", shouldUpdate: false},
		//{name: "NestedDeferWithLockAndUnlock", testPath: "e2e-tests/defer/NestedDeferWithLockAndUnlock/prog1.go", resPath: "e2e-tests/defer/NestedDeferWithLockAndUnlock/prog1_expected.json", shouldUpdate: false},
		//{name: "NestedDeferWithLockAndUnlockAndGoroutine", testPath: "e2e-tests/defer/NestedDeferWithLockAndUnlockAndGoroutine/prog1.go", resPath: "e2e-tests/defer/NestedDeferWithLockAndUnlockAndGoroutine/prog1_expected.json", shouldUpdate: false},
		//{name: "LockAndUnlockIfMap", testPath: "e2e-tests/locksAndUnlocks/LockAndUnlockIfMap/prog1.go", resPath: "e2e-tests/locksAndUnlocks/LockAndUnlockIfMap/prog1_expected.json", shouldUpdate: false},
		//{name: "MultipleLocksNoRace", testPath: "e2e-tests/locksAndUnlocks/MultipleLocksNoRace/prog1.go", resPath: "e2e-tests/locksAndUnlocks/MultipleLocksNoRace/prog1_expected.json", shouldUpdate: false},
		//{name: "MultipleLocksRace", testPath: "e2e-tests/locksAndUnlocks/MultipleLocksRace/prog1.go", resPath: "e2e-tests/locksAndUnlocks/MultipleLocksRace/prog1_expected.json", shouldUpdate: false	},
		//{name: "NestedFunctions", testPath: "e2e-tests/general/NestedFunctions/prog1.go", resPath: "e2e-tests/general/NestedFunctions/prog1_expected.json", shouldUpdate: false},
		//{name: "NestedConditionWithLockInAllBranches", testPath: "e2e-tests/locksAndUnlocks/NestedConditionWithLockInAllBranches/prog1.go", resPath: "e2e-tests/locksAndUnlocks/NestedConditionWithLockInAllBranches/prog1_expected.json", shouldUpdate: false},
		//{name: "DataRaceMap", testPath: "e2e-tests/general/DataRaceMap/prog1.go", resPath: "e2e-tests/general/DataRaceMap/prog1_expected.json", shouldUpdate: false},
		//{name: "ForLoop", testPath: "e2e-tests/ForLoop/prog1.go", resPath: "e2e-tests/ForLoop/prog1_expected.json", shouldUpdate: false},
		//{name: "DataRaceShadowedErr", testPath: "e2e-tests/general/DataRaceShadowedErr/prog1.go", resPath: "e2e-tests/general/DataRaceShadowedErr/prog1_expected.json", shouldUpdate: false},
		//{name: "DataRaceInterfaceOverChannel", testPath: "e2e-tests/pointerAnalysis/DataRaceInterfaceOverChannel/prog1.go", resPath: "e2e-tests/pointerAnalysis/DataRaceInterfaceOverChannel/prog1_expected.json", shouldUpdate: true},
		//{name: "DataRaceProperty", testPath: "e2e-tests/general/DataRaceProperty/prog1.go", resPath: "e2e-tests/general/DataRaceProperty/prog1_expected.json", shouldUpdate: false},
		//{name: "DataRaceWithOnlyAlloc", testPath: "e2e-tests/general/DataRaceWithOnlyAlloc/prog1.go", resPath: "e2e-tests/general/DataRaceWithOnlyAlloc/prog1_expected.json", shouldUpdate: false},
		//{name: "LockInsideGoroutine", testPath: "e2e-tests/locksAndUnlocks/LockInsideGoroutine/prog1.go", resPath: "e2e-tests/locksAndUnlocks/LockInsideGoroutine/prog1_expected.json", shouldUpdate: false},
		//{name: "DataRaceWithSameFunction", testPath: "e2e-tests/general/DataRaceWithSameFunction/prog1.go", resPath: "e2e-tests/general/DataRaceWithSameFunction/prog1_expected.json", shouldUpdate: false},
		//{name: "StructMethod", testPath: "e2e-tests/general/StructMethod/prog1.go", resPath: "e2e-tests/general/StructMethod/prog1_expected.json", shouldUpdate: false},
		//{name: "DataRaceIceCreamMaker", testPath: "e2e-tests/interfaces/DataRaceIceCreamMaker/prog1.go", resPath: "e2e-tests/interfaces/DataRaceIceCreamMaker/prog1_expected.json", shouldUpdate: false},
		//{name: "InterfaceWithLock", testPath: "e2e-tests/interfaces/InterfaceWithLock/prog1.go", resPath: "e2e-tests/interfaces/InterfaceWithLock/prog1_expected.json", shouldUpdate: false},
		//{name: "NestedInterface", testPath: "e2e-tests/interfaces/NestedInterface/prog1.go", resPath: "e2e-tests/interfaces/NestedInterface/prog1_expected.json", shouldUpdate: false},
		//{name: "TestNoRaceStackPushPop", testPath: "e2e-tests/stdlib/TestNoRaceStackPushPop/prog1.go", resPath: "e2e-tests/stdlib/TestNoRaceStackPushPop/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceComplement", testPath: "e2e-tests/stdlib/TestRaceComplement/prog1.go", resPath: "e2e-tests/stdlib/TestRaceComplement/prog1_expected.json", shouldUpdate: true},
		//{name: "RaceNestedArrayCopy", testPath: "e2e-tests/stdlib/RaceNestedArrayCopy/prog1.go", resPath: "e2e-tests/stdlib/RaceNestedArrayCopy/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceNestedStruct", testPath: "e2e-tests/stdlib/TestRaceNestedStruct/prog1.go", resPath: "e2e-tests/stdlib/TestRaceNestedStruct/prog1_expected.json", shouldUpdate: true},
		//{name: "TestNoRaceAsFunc4", testPath: "e2e-tests/stdlib/TestNoRaceAsFunc4/prog1.go", resPath: "e2e-tests/stdlib/TestNoRaceAsFunc4/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceAsFunc3", testPath: "e2e-tests/stdlib/TestRaceAsFunc3/prog1.go", resPath: "e2e-tests/stdlib/TestRaceAsFunc3/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceAsFunc2", testPath: "e2e-tests/stdlib/TestRaceAsFunc2/prog1.go", resPath: "e2e-tests/stdlib/TestRaceAsFunc2/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceAsFunc1", testPath: "e2e-tests/stdlib/TestRaceAsFunc1/prog1.go", resPath: "e2e-tests/stdlib/TestRaceAsFunc1/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceCaseTypeIssue5890", testPath: "e2e-tests/stdlib/TestRaceCaseTypeIssue5890/prog1.go", resPath: "e2e-tests/stdlib/TestRaceCaseTypeIssue5890/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceCaseIssue6418", testPath: "e2e-tests/stdlib/TestRaceCaseIssue6418/prog1.go", resPath: "e2e-tests/stdlib/TestRaceCaseIssue6418/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceCaseFallthrough", testPath: "e2e-tests/stdlib/TestRaceCaseFallthrough/prog1.go", resPath: "e2e-tests/stdlib/TestRaceCaseFallthrough/prog1_expected.json", shouldUpdate: true},
		//{name: "TestNoRaceBlank", testPath: "e2e-tests/unsupported/TestNoRaceBlank/prog1.go", resPath: "e2e-tests/unsupported/TestNoRaceBlank/prog1_expected.json", shouldUpdate: true}, // blank space
		//{name: "TestRaceInterCall", testPath: "e2e-tests/stdlib/TestRaceInterCall/prog1.go", resPath: "e2e-tests/stdlib/TestRaceInterCall/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceMethodThunk4", testPath: "e2e-tests/stdlib/TestRaceMethodThunk4/prog1.go", resPath: "e2e-tests/stdlib/TestRaceMethodThunk4/prog1_expected.json", shouldUpdate: true}, // Might be a bug in pointer analysis
		//{name: "TestRaceMethodThunk3", testPath: "e2e-tests/stdlib/TestRaceMethodThunk3/prog1.go", resPath: "e2e-tests/stdlib/TestRaceMethodThunk3/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceMethodThunk2", testPath: "e2e-tests/stdlib/TestRaceMethodThunk2/prog1.go", resPath: "e2e-tests/stdlib/TestRaceMethodThunk4/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceMethodThunk", testPath: "e2e-tests/stdlib/TestRaceMethodThunk/prog1.go", resPath: "e2e-tests/stdlib/TestRaceMethodThunk/prog1_expected.json", shouldUpdate: true},
		//{name: "TestNoRaceMethodThunk", testPath: "e2e-tests/stdlib/TestNoRaceMethodThunk/prog1.go", resPath: "e2e-tests/stdlib/TestNoRaceMethodThunk/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceNestedStruct", testPath: "e2e-tests/stdlib/TestRaceNestedStruct/prog1.go", resPath: "e2e-tests/stdlib/TestRaceNestedStruct/prog1_expected.json", shouldUpdate: true},
		//{name: "TestNoRaceEmptyStruct", testPath: "e2e-tests/stdlibNoSuccess/TestNoRaceEmptyStruct/prog1.go", resPath: "e2e-tests/stdlibNoSuccess/TestNoRaceEmptyStruct/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceHeapParam", testPath: "e2e-tests/stdlib/TestRaceHeapParam/prog1.go", resPath: "e2e-tests/stdlib/TestRaceHeapParam/prog1_expected.json", shouldUpdate: true}, // No ssa param as value
		//{name: "TestRaceStructInd", testPath: "e2e-tests/stdlib/TestRaceStructInd/prog1.go", resPath: "e2e-tests/stdlib/TestRaceStructInd/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceAppendSliceStruct", testPath: "e2e-tests/stdlibNoSuccess/TestRaceAppendSliceStruct/prog1.go", resPath: "e2e-tests/stdlibNoSuccess/TestRaceAppendSliceStruct/prog1_expected.json", shouldUpdate: true}, // spread operator can't tell which item are affected
		//{name: "TestRaceSliceStruct", testPath: "e2e-tests/stdlibNoSuccess/TestRaceSliceStruct/prog1.go", resPath: "e2e-tests/stdlibNoSuccess/TestRaceSliceStruct/prog1_expected.json", shouldUpdate: true}, // same
		//{name: "TestRaceSliceString", testPath: "e2e-tests/stdlib/TestRaceSliceString/prog1.go", resPath: "e2e-tests/stdlib/TestRaceSliceString/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceSliceSlice2", testPath: "e2e-tests/stdlib/TestRaceSliceSlice2/prog1.go", resPath: "e2e-tests/stdlib/TestRaceSliceSlice2/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceSliceSlice", testPath: "e2e-tests/stdlib/TestRaceSliceSlice/prog1.go", resPath: "e2e-tests/stdlib/TestRaceSliceSlice/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceBlockAs", testPath: "e2e-tests/stdlib/TestRaceBlockAs/prog1.go", resPath: "e2e-tests/stdlib/TestRaceBlockAs/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceTypeAssert", testPath: "e2e-tests/stdlib/TestRaceTypeAssert/prog1.go", resPath: "e2e-tests/stdlib/TestRaceTypeAssert/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceAddrExpr", testPath: "e2e-tests/stdlib/TestRaceAddrExpr/prog1.go", resPath: "e2e-tests/stdlib/TestRaceAddrExpr/prog1_expected.json", shouldUpdate: true},
		//{name: "TestNoRaceAddrExpr", testPath: "e2e-tests/stdlib/TestNoRaceAddrExpr/prog1.go", resPath: "e2e-tests/stdlib/TestNoRaceAddrExpr/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceDeferArg2", testPath: "e2e-tests/stdlib/TestRaceDeferArg2/prog1.go", resPath: "e2e-tests/stdlib/TestRaceDeferArg2/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceDeferArg", testPath: "e2e-tests/stdlib/TestRaceDeferArg/prog1.go", resPath: "e2e-tests/stdlib/TestRaceDeferArg/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRacePanicArg", testPath: "e2e-tests/stdlib/TestRacePanicArg/prog1.go", resPath: "e2e-tests/stdlib/TestRacePanicArg/prog1_expected.json", shouldUpdate: true},
		//{name: "TestNoRaceMethodValue", testPath: "e2e-tests/stdlib/TestNoRaceMethodValue/prog1.go", resPath: "e2e-tests/stdlib/TestNoRaceMethodValue/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceMethodValue3", testPath: "e2e-tests/stdlib/TestRaceMethodValue3/prog1.go", resPath: "e2e-tests/stdlib/TestRaceMethodValue3/prog1_expected.json", shouldUpdate: true}, // // Might be a bug in pointer analysis
		//{name: "TestRaceMethodValue2", testPath: "e2e-tests/stdlib/TestRaceMethodValue2/prog1.go", resPath: "e2e-tests/stdlib/TestRaceMethodValue2/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceMethodValue", testPath: "e2e-tests/stdlib/TestRaceMethodValue/prog1.go", resPath: "e2e-tests/stdlib/TestRaceMethodValue/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceMethodCall2", testPath: "e2e-tests/stdlib/TestRaceMethodCall2/prog1.go", resPath: "e2e-tests/stdlib/TestRaceMethodCall2/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceMethodCall", testPath: "e2e-tests/stdlib/TestRaceMethodCall/prog1.go", resPath: "e2e-tests/stdlib/TestRaceMethodCall/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceFuncCall", testPath: "e2e-tests/stdlib/TestRaceFuncCall/prog1.go", resPath: "e2e-tests/stdlib/TestRaceFuncCall/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceInterCall2", testPath: "e2e-tests/stdlib/TestRaceInterCall2/prog1.go", resPath: "e2e-tests/stdlib/TestRaceInterCall2/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceInterCall", testPath: "e2e-tests/stdlib/TestRaceInterCall/prog1.go", resPath: "e2e-tests/stdlib/TestRaceInterCall/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceMapInit2", testPath: "e2e-tests/stdlib/TestRaceMapInit2/prog1.go", resPath: "e2e-tests/stdlib/TestRaceMapInit2/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceMapInit", testPath: "e2e-tests/stdlib/TestRaceMapInit/prog1.go", resPath: "e2e-tests/stdlib/TestRaceMapInit/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceArrayInit", testPath: "e2e-tests/stdlib/TestRaceArrayInit/prog1.go", resPath: "e2e-tests/stdlib/TestRaceArrayInit/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceStructInit", testPath: "e2e-tests/stdlib/TestRaceStructInit/prog1.go", resPath: "e2e-tests/stdlib/TestRaceStructInit/prog1_expected.json", shouldUpdate: true},
		//{name: "TestNoRaceFuncUnlock", testPath: "e2e-tests/stdlibNoSuccess/TestNoRaceFuncUnlock/prog1.go", resPath: "e2e-tests/stdlibNoSuccess/TestNoRaceFuncUnlock/prog1_expected.json", shouldUpdate: true}, // No pointer analysis for locks
		//{name: "TestRaceFuncItself", testPath: "e2e-tests/stdlib/TestRaceFuncItself/prog1.go", resPath: "e2e-tests/stdlib/TestRaceFuncItself/prog1_expected.json", shouldUpdate: true},
		//{name: "TestNoRaceShortCalc2", testPath: "e2e-tests/stdlibNoSuccess/TestNoRaceShortCalc2/prog1.go", resPath: "e2e-tests/stdlibNoSuccess/TestNoRaceShortCalc2/prog1_expected.json", shouldUpdate: true}, // Cant evaluate first part of condition to see the second will never execute
		//{name: "TestNoRaceShortCalc", testPath: "e2e-tests/stdlibNoSuccess/TestNoRaceShortCalc/prog1.go", resPath: "e2e-tests/stdlibNoSuccess/TestNoRaceShortCalc/prog1_expected.json", shouldUpdate: true}, // Cant evaluate first part of condition to see the second will never execute
		//{name: "TestNoRaceOr", testPath: "e2e-tests/stdlibNoSuccess/TestNoRaceOr/prog1.go", resPath: "e2e-tests/stdlibNoSuccess/TestNoRaceOr/prog1_expected.json", shouldUpdate: true}, // Cant evaluate first part of condition to see the second will never execute
		//{name: "TestRaceOr2", testPath: "e2e-tests/stdlib/TestRaceOr2/prog1.go", resPath: "e2e-tests/stdlib/TestRaceOr2/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceOr", testPath: "e2e-tests/stdlib/TestRaceOr/prog1.go", resPath: "e2e-tests/stdlib/TestRaceOr/prog1_expected.json", shouldUpdate: true},
		//{name: "TestNoRaceAnd", testPath: "e2e-tests/stdlibNoSuccess/TestNoRaceAnd/prog1.go", resPath: "e2e-tests/stdlibNoSuccess/TestNoRaceAnd/prog1_expected.json", shouldUpdate: true}, // Cant evaluate first part of condition to see the second will never execute
		//{name: "TestRaceAnd2", testPath: "e2e-tests/stdlib/TestNoRaceShortCalc/prog1.go", resPath: "e2e-tests/stdlib/TestNoRaceShortCalc/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceAnd", testPath: "e2e-tests/stdlib/TestRaceAnd/prog1.go", resPath: "e2e-tests/stdlib/TestRaceAnd/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceEmptyInterface2", testPath: "e2e-tests/stdlibNoSuccess/TestRaceEmptyInterface2/prog1.go", resPath: "e2e-tests/stdlibNoSuccess/TestRaceEmptyInterface2/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceEmptyInterface1", testPath: "e2e-tests/stdlib/TestRaceEmptyInterface/prog1.go", resPath: "e2e-tests/stdlib/TestRaceEmptyInterface/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceRune", testPath: "e2e-tests/stdlib/TestRaceRune/prog1.go", resPath: "e2e-tests/stdlib/TestRaceRune/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceIndirection", testPath: "e2e-tests/stdlibNoSuccess/TestRaceIndirection/prog1.go", resPath: "e2e-tests/stdlibNoSuccess/TestRaceIndirection/prog1_expected.json", shouldUpdate: true}, // sync using channels
		//{name: "TestRaceFuncArgsRW", testPath: "e2e-tests/stdlib/TestRaceFuncArgsRW/prog1.go", resPath: "e2e-tests/stdlib/TestRaceFuncArgsRW/prog1_expected.json", shouldUpdate: true},
		//{name: "TestNoRaceFuncArgsRW", testPath: "e2e-tests/stdlib/TestNoRaceFuncArgsRW/prog1.go", resPath: "e2e-tests/stdlib/TestNoRaceFuncArgsRW/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceAppendCapRW", testPath: "e2e-tests/stdlib/TestRaceAppendCapRW/prog1.go", resPath: "e2e-tests/stdlib/TestRaceAppendCapRW/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceAppendLenRW", testPath: "e2e-tests/stdlib/TestRaceAppendLenRW/prog1.go", resPath: "e2e-tests/stdlib/TestRaceAppendLenRW/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceAppendRW", testPath: "e2e-tests/stdlib/TestRaceAppendRW/prog1.go", resPath: "e2e-tests/stdlib/TestRaceAppendRW/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRacePanic", testPath: "e2e-tests/stdlib/TestRacePanic/prog1.go", resPath: "e2e-tests/stdlib/TestRacePanic/prog1_expected.json", shouldUpdate: true}, // cfg is weird because of the recover
		//{name: "TestRaceFuncVariableWW", testPath: "e2e-tests/stdlib/TestRaceFuncVariableWW/prog1.go", resPath: "e2e-tests/stdlib/TestRaceFuncVariableWW/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceFuncVariableRW", testPath: "e2e-tests/stdlib/TestRaceFuncVariableRW/prog1.go", resPath: "e2e-tests/stdlib/TestRaceFuncVariableRW/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceUnsafePtrRW", testPath: "e2e-tests/stdlib/TestRaceUnsafePtrRW/prog1.go", resPath: "e2e-tests/stdlib/TestRaceUnsafePtrRW/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceComplex128WW", testPath: "e2e-tests/stdlib/TestRaceComplex128WW/prog1.go", resPath: "e2e-tests/stdlib/TestRaceComplex128WW/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceFloat64WW", testPath: "e2e-tests/stdlib/TestRaceFloat64WW/prog1.go", resPath: "e2e-tests/stdlib/TestRaceFloat64WW/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceStringPtrRW", testPath: "e2e-tests/stdlib/TestRaceStringPtrRW/prog1.go", resPath: "e2e-tests/stdlib/TestRaceStringPtrRW/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceStringRW", testPath: "e2e-tests/stdlib/TestRaceStringRW/prog1.go", resPath: "e2e-tests/stdlib/TestRaceStringRW/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceIntptrRW", testPath: "e2e-tests/stdlib/TestRaceIntptrRW/prog1.go", resPath: "e2e-tests/stdlib/TestRaceIntptrRW/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceError", testPath: "e2e-tests/stdlib/TestRaceError/prog1.go", resPath: "e2e-tests/stdlib/TestRaceError/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceIfaceConv", testPath: "e2e-tests/stdlib/TestRaceIfaceConv/prog1.go", resPath: "e2e-tests/stdlib/TestRaceIfaceConv/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceEfaceConv", testPath: "e2e-tests/stdlib/TestRaceEfaceConv/prog1.go", resPath: "e2e-tests/stdlib/TestRaceEfaceConv/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceIfaceCmpNil", testPath: "e2e-tests/stdlib/TestRaceIfaceCmpNil/prog1.go", resPath: "e2e-tests/stdlib/TestRaceIfaceCmpNil/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceIfaceCmp", testPath: "e2e-tests/stdlib/TestRaceIfaceCmp/prog1.go", resPath: "e2e-tests/stdlib/TestRaceIfaceCmp/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceIfaceWW", testPath: "e2e-tests/stdlib/TestRaceIfaceWW/prog1.go", resPath: "e2e-tests/stdlib/TestRaceIfaceWW/prog1_expected.json", shouldUpdate: true}, // Before write, a read is performed. So the creation confused with the read later.
		//{name: "TestRaceEfaceWW", testPath: "e2e-tests/stdlib/TestRaceEfaceWW/prog1.go", resPath: "e2e-tests/stdlib/TestRaceEfaceWW/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceStructFieldRW3", testPath: "e2e-tests/stdlib/TestRaceStructFieldRW3/prog1.go", resPath: "e2e-tests/stdlib/TestRaceStructFieldRW3/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceStructFieldRW2", testPath: "e2e-tests/stdlib/TestRaceStructFieldRW2/prog1.go", resPath: "e2e-tests/stdlib/TestRaceStructFieldRW2/prog1_expected.json", shouldUpdate: true},
		//{name: "TestNoRaceStructFieldRW2", testPath: "e2e-tests/stdlibNoSuccess/TestNoRaceStructFieldRW2/prog1.go", resPath: "e2e-tests/stdlibNoSuccess/TestNoRaceStructFieldRW2/prog1_expected.json", shouldUpdate: true}, // Not sure. Probably because it's dereference by value an not pointer
		//{name: "TestNoRaceStructFieldRW1", testPath: "e2e-tests/stdlibNoSuccess/TestNoRaceStructFieldRW1/prog1.go", resPath: "e2e-tests/stdlibNoSuccess/TestNoRaceStructFieldRW1/prog1_expected.json", shouldUpdate: true}, // Not sure. Probably because it's dereference by value an not pointer
		//{name: "TestRaceStructFieldRW1", testPath: "e2e-tests/stdlib/TestRaceStructFieldRW1/prog1.go", resPath: "e2e-tests/stdlib/TestRaceStructFieldRW1/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceStructRW", testPath: "e2e-tests/stdlib/TestRaceStructRW/prog1.go", resPath: "e2e-tests/stdlib/TestRaceStructRW/prog1_expected.json", shouldUpdate: true}, // Due to pointer analysis struct init part are counted as struct itself
		//{name: "TestRaceArrayCopy", testPath: "e2e-tests/stdlib/TestRaceArrayCopy/prog1.go", resPath: "e2e-tests/stdlib/TestRaceArrayCopy/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceSprint", testPath: "e2e-tests/stdlib/TestRaceSprint/prog1.go", resPath: "e2e-tests/stdlib/TestRaceSprint/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceFuncArgument2", testPath: "e2e-tests/stdlib/TestRaceFuncArgument2/prog1.go", resPath: "e2e-tests/stdlib/TestRaceFuncArgument2/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceFuncArgument", testPath: "e2e-tests/stdlib/TestRaceFuncArgument/prog1.go", resPath: "e2e-tests/stdlib/TestRaceFuncArgument/prog1_expected.json", shouldUpdate: true},
		//{name: "TestNoRaceEnoughRegisters", testPath: "e2e-tests/stdlib/TestNoRaceEnoughRegisters/prog1.go", resPath: "e2e-tests/stdlib/TestNoRaceEnoughRegisters/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceRotate", testPath: "e2e-tests/stdlib/TestRaceRotate/prog1.go", resPath: "e2e-tests/stdlib/TestRaceRotate/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceModConst", testPath: "e2e-tests/stdlib/TestRaceModConst/prog1.go", resPath: "e2e-tests/stdlib/TestRaceModConst/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceMod", testPath: "e2e-tests/stdlib/TestRaceMod/prog1.go", resPath: "e2e-tests/stdlib/TestRaceMod/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceDivConst", testPath: "e2e-tests/stdlib/TestRaceDivConst/prog1.go", resPath: "e2e-tests/stdlib/TestRaceDivConst/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceDiv", testPath: "e2e-tests/stdlib/TestRaceDiv/prog1.go", resPath: "e2e-tests/stdlib/TestRaceDiv/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRaceComplement", testPath: "e2e-tests/stdlib/TestRaceComplement/prog1.go", resPath: "e2e-tests/stdlib/TestRaceComplement/prog1_expected.json", shouldUpdate: true},
		//{name: "TestNoRacePlus", testPath: "e2e-tests/stdlib/TestNoRacePlus/prog1.go", resPath: "e2e-tests/stdlib/TestNoRacePlus/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRacePlus2", testPath: "e2e-tests/stdlib/TestRacePlus2/prog1.go", resPath: "e2e-tests/stdlib/TestRacePlus2/prog1_expected.json", shouldUpdate: true},
		//{name: "TestRacePlus", testPath: "e2e-tests/stdlib/TestRacePlus/prog1.go", resPath: "e2e-tests/stdlib/TestRacePlus/prog1_expected.json", shouldUpdate: true},
		{name: "TestRaceCaseTypeBody", testPath: "e2e-tests/stdlib/TestRaceCaseTypeBody/prog1.go", resPath: "e2e-tests/stdlib/TestRaceCaseTypeBody/prog1_expected.json", shouldUpdate: true},
		{name: "TestRaceCaseType", testPath: "e2e-tests/stdlib/TestRaceCaseType/prog1.go", resPath: "e2e-tests/stdlib/TestRaceCaseType/prog1_expected.json", shouldUpdate: true},
		//{name: "TestNoRaceCaseFallthrough", testPath: "e2e-tests/stdlib/stdlibNoSuccess/prog1.go", resPath: "e2e-tests/stdlibNoSuccess/TestNoRaceCaseFallthrough/prog1_expected.json", shouldUpdate: true}, // No way to determine flow as the detector is flow insensitive
		{name: "TestRaceCaseBody", testPath: "e2e-tests/stdlib/TestRaceCaseBody/prog1.go", resPath: "e2e-tests/stdlib/TestRaceCaseBody/prog1_expected.json", shouldUpdate: true},
		{name: "TestRaceCaseCondition2", testPath: "e2e-tests/stdlib/TestRaceCaseCondition2/prog1.go", resPath: "e2e-tests/stdlib/TestRaceCaseCondition2/prog1_expected.json", shouldUpdate: true},
		{name: "TestRaceCaseCondition", testPath: "e2e-tests/stdlib/TestRaceCaseCondition/prog1.go", resPath: "e2e-tests/stdlib/TestRaceCaseCondition/prog1_expected.json", shouldUpdate: true},
		{name: "TestRaceInt32RWClosures", testPath: "e2e-tests/stdlib/TestRaceInt32RWClosures/prog1.go", resPath: "e2e-tests/stdlib/TestRaceInt32RWClosures/prog1_expected.json", shouldUpdate: true},
		{name: "TestNoRaceIntRWClosures", testPath: "e2e-tests/stdlib/TestNoRaceIntRWClosures/prog1.go", resPath: "e2e-tests/stdlib/TestNoRaceIntRWClosures/prog1_expected.json", shouldUpdate: true},
		{name: "TestRaceIntRWClosures", testPath: "e2e-tests/stdlib/TestRaceIntRWClosures/prog1.go", resPath: "e2e-tests/stdlib/TestRaceIntRWClosures/prog1_expected.json", shouldUpdate: true},
		{name: "TestRaceIntRWGlobalFuncs", testPath: "e2e-tests/stdlib/TestRaceIntRWGlobalFuncs/prog1.go", resPath: "e2e-tests/stdlib/TestRaceIntRWGlobalFuncs/prog1_expected.json", shouldUpdate: true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			domain.GoroutineCounter.Reset()
			domain.GuardedAccessCounter.Reset()

			ssaProg, ssaPkg, err := ssaUtils.LoadPackage(tc.testPath)
			require.NoError(t, err)
			ssaUtils.SetGlobalProgram(ssaProg)

			entryFunc := ssaPkg.Func("main")
			entryCallCommon := ssa.CallCommon{Value: entryFunc}
			functionState := ssaUtils.HandleCallCommon(domain.NewEmptyGoroutineState(), &entryCallCommon)
			testresult := testutils.TestResult{Lockset: functionState.Lockset, GuardedAccess: functionState.GuardedAccesses}
			dump, err := json.MarshalIndent(testresult, "", "\t")
			require.NoError(t, err)
			if shouldUpdateAll == updateAll || shouldUpdateAll == testSelective && tc.shouldUpdate {
				utils.UpdateFile(t, tc.resPath, dump)
			}
			testutils.CompareResult(t, tc.resPath, functionState.Lockset, functionState.GuardedAccesses)
			Analysis(ssaPkg, ssaProg, functionState.GuardedAccesses)
		})
	}
}
