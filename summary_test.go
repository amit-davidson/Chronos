package main

import (
	"StaticRaceDetector/domain"
	"StaticRaceDetector/testutils"
	"StaticRaceDetector/utils"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

var shouldUpdate = true

func TestGetFunctionSummary(t *testing.T) {
	var testCases = []struct {
		name     string
		testPath string
		resPath  string
	}{
		//{name: "Lock", testPath: "testutils/Lock/prog1.go", resPath: "testutils/Lock/prog1_expected"},
		//{name: "LockAndUnlock", testPath: "testutils/LockAndUnlock/prog1.go", resPath: "testutils/LockAndUnlock/prog1_expected"},
		//{name: "LockAndUnlockIfBranch", testPath: "testutils/LockAndUnlockIfBranch/prog1.go", resPath: "testutils/LockAndUnlockIfBranch/prog1_expected"},
		//{name: "LockAndUnlockIfMap", testPath: "testutils/LockAndUnlockIfMap/prog1.go", resPath: "testutils/LockAndUnlockIfMap/prog1_expected"},
		//{name: "NestedFunctions", testPath: "testutils/NestedFunctions/prog1.go", resPath: "testutils/NestedFunctions/prog1_expected"},
		{name: "DataRaceMap", testPath: "testutils/DataRaceMap/prog1.go", resPath: "testutils/DataRaceMap/prog1_expected"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ssaProg, ssaPkg, err := utils.LoadPackage(tc.testPath)
			require.NoError(t, err)

			entryFunc := ssaPkg.Func("main")
			emptyLS := domain.NewEmptyLockSet()
			lsRet, guardedAccessRet := GetFunctionSummary(entryFunc, emptyLS, utils.GetUUID())

			testresult := testutils.TestResult{Lockset: lsRet, GuardedAccess: guardedAccessRet}
			dump, err := json.Marshal(testresult)
			require.NoError(t, err)
			if shouldUpdate {
				utils.UpdateFile(t, tc.resPath, dump)
			}
			testutils.CompareResult(t, tc.resPath, lsRet, guardedAccessRet)
			Analysis(ssaPkg, ssaProg, guardedAccessRet)
		})
	}
}
