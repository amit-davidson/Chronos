package main

import (
	"StaticRaceDetector/domain"
	"StaticRaceDetector/testutils"
	"StaticRaceDetector/utils"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"go/token"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
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
		{name: "LockAndUnlockIfMap", testPath: "testutils/LockAndUnlockIfMap/prog1.go", resPath: "testutils/LockAndUnlockIfMap/prog1_expected"},
		//{name: "NestedFunctions", testPath: "testutils/NestedFunctions/prog1.go", resPath: "testutils/NestedFunctions/prog1_expected"},
		//{name: "NestedFunctionsTest", testPath: "testutils/NestedFunctionsTest/prog1.go", resPath: "testutils/NestedFunctionsTest/prog1_expected"},
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

			config := &pointer.Config{
				Mains: []*ssa.Package{ssaPkg},
			}

			valuesQueriesToGuardAccess := map[token.Pos]*domain.GuardedAccess{}
			for _, guardedAccess := range guardedAccessRet {
				if pointer.CanPoint(guardedAccess.Value.Type()) {
					config.AddQuery(guardedAccess.Value)
					valuesQueriesToGuardAccess[guardedAccess.Value.Pos()] = guardedAccess
				}
			}

			positionsToGuardAccesses := map[token.Pos][]*domain.GuardedAccess{}
			result, err := pointer.Analyze(config)
			if err != nil {
				panic(err) // internal error in pointer analysis
			}

			for v, l := range result.Queries {
				for _, label := range l.PointsTo().Labels() {
					guardedAccess := valuesQueriesToGuardAccess[v.Pos()]
					allocPos := label.Value()
					positionsToGuardAccesses[allocPos.Pos()] = append(positionsToGuardAccesses[allocPos.Pos()], guardedAccess)
				}
			}
			for _, guardedAccesses := range positionsToGuardAccesses {
				for _, guardedAccessesA := range guardedAccesses {
					for _, guardedAccessesB := range guardedAccesses {
						if !guardedAccessesA.Intersects(guardedAccessesB) {
							valueA := guardedAccessesA.Value
							valueB := guardedAccessesB.Value
							label := fmt.Sprintf(" %s with pos:%s has race condition with %s pos:%s \n", valueA, ssaProg.Fset.Position(valueA.Pos()), valueB, ssaProg.Fset.Position(valueB.Pos()))
							print(label)
						}
					}
				}
			}
		})
	}
}
