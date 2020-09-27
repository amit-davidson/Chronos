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
		//{name: "Lock", testPath: "testutils/Lock/prog1.go", resPath: "testutils/Lock/prog1_expected.json", shouldUpdate: false},
		//{name: "LockAndUnlock", testPath: "testutils/LockAndUnlock/prog1.go", resPath: "testutils/LockAndUnlock/prog1_expected.json", shouldUpdate: false},
		//{name: "LockAndUnlockIfBranch", testPath: "testutils/LockAndUnlockIfBranch/prog1.go", resPath: "testutils/LockAndUnlockIfBranch/prog1_expected.json", shouldUpdate: false},
		//{name: "ElseIf", testPath: "testutils/ElseIf/prog1.go", resPath: "testutils/ElseIf/prog1_expected.json", shouldUpdate: false},
		//{name: "DeferredLockAndUnlockIfBranch", testPath: "testutils/DeferredLockAndUnlockIfBranch/prog1.go", resPath: "testutils/DeferredLockAndUnlockIfBranch/prog1_expected.json", shouldUpdate: false}, // Not tested
		//{name: "NestedDeferWithLockAndUnlock", testPath: "testutils/NestedDeferWithLockAndUnlock/prog1.go", resPath: "testutils/NestedDeferWithLockAndUnlock/prog1_expected.json", shouldUpdate: false}, // Not tested
		//{name: "NestedDeferWithLockAndUnlockAndGoroutine", testPath: "testutils/NestedDeferWithLockAndUnlockAndGoroutine/prog1.go", resPath: "testutils/NestedDeferWithLockAndUnlockAndGoroutine/prog1_expected.json", shouldUpdate: false}, // Not tested
		//{name: "LockAndUnlockIfMap", testPath: "testutils/LockAndUnlockIfMap/prog1.go", resPath: "testutils/LockAndUnlockIfMap/prog1_expected.json", shouldUpdate: false},
		//{name: "NestedFunctions", testPath: "testutils/NestedFunctions/prog1.go", resPath: "testutils/NestedFunctions/prog1_expected.json", shouldUpdate: false},
		{name: "DataRaceMap", testPath: "testutils/DataRaceMap/prog1.go", resPath: "testutils/DataRaceMap/prog1_expected.json", shouldUpdate: true},
		//{name: "ForLoop", testPath: "testutils/ForLoop/prog1.go", resPath: "testutils/ForLoop/prog1_expected.json", shouldUpdate: true},
		//{name: "DataRaceShadowedErr", testPath: "testutils/DataRaceShadowedErr/prog1.go", resPath: "testutils/DataRaceShadowedErr/prog1_expected.json", shouldUpdate: false},
		//{name: "DataRaceProperty", testPath: "testutils/DataRaceProperty/prog1.go", resPath: "testutils/DataRaceProperty/prog1_expected.json", shouldUpdate: false},
		//{name: "DataRaceWithOnlyAlloc", testPath: "testutils/DataRaceWithOnlyAlloc/prog1.go", resPath: "testutils/DataRaceWithOnlyAlloc/prog1_expected.json", shouldUpdate: false},
		//{name: "DataRaceWithSameFunction", testPath: "testutils/DataRaceWithSameFunction/prog1.go", resPath: "testutils/DataRaceWithSameFunction/prog1_expected.json", shouldUpdate: false},
		//{name: "StructMethod", testPath: "testutils/StructMethod/prog1.go", resPath: "testutils/StructMethod/prog1_expected.json", shouldUpdate: false},
		//{name: "DataRaceIceCreamMaker", testPath: "testutils/DataRaceIceCreamMaker/prog1.go", resPath: "testutils/DataRaceIceCreamMaker/prog1_expected.json", shouldUpdate: false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			domain.GoroutineCounter.Reset()
			domain.GuardedAccessCounter.Reset()

			ssaProg, ssaPkg, err := ssaUtils.LoadPackage(tc.testPath)
			_ = ssaProg
			require.NoError(t, err)

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
			//Analysis(ssaPkg, ssaProg, guardedAccesses)
		})
	}
}
