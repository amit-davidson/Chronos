package main

import (
	"StaticRaceDetector/domain"
	"StaticRaceDetector/testutils"
	"StaticRaceDetector/utils"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"go/token"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"testing"
)

var shouldUpdate = true

func TestGetFunctionSummary(t *testing.T) {
	var testCases = []struct {
		name     string
		testPath string
		resPath  string
	}{
		{name: "Lock", testPath: "testutils/Lock/prog1.go", resPath: "testutils/Lock/prog1_expected"},
		{name: "LockAndUnlock", testPath: "testutils/LockAndUnlock/prog1.go", resPath: "testutils/LockAndUnlock/prog1_expected"},
		{name: "LockAndUnlockIfBranch", testPath: "testutils/LockAndUnlockIfBranch/prog1.go", resPath: "testutils/LockAndUnlockIfBranch/prog1_expected"},
		{name: "LockAndUnlockIfMap", testPath: "testutils/LockAndUnlockIfMap/prog1.go", resPath: "testutils/LockAndUnlockIfMap/prog1_expected"},
		{name: "NestedFunctions", testPath: "testutils/NestedFunctions/prog1.go", resPath: "testutils/NestedFunctions/prog1_expected"},
		{name: "NestedFunctionsTest", testPath: "testutils/NestedFunctionsTest/prog1.go", resPath: "testutils/NestedFunctionsTest/prog1_expected"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var conf loader.Config

			file, err := conf.ParseFile(tc.testPath, nil)
			if err != nil {
				fmt.Print(err) // parse error
				return
			}

			// Create single-file main package and import its dependencies.
			conf.CreateFromFiles("main", file)

			iprog, err := conf.Load()
			if err != nil {
				fmt.Print(err) // type error in some package
				return
			}

			// Create SSA-form program representation.
			prog := ssautil.CreateProgram(iprog, 0)
			ssaPkg := prog.Package(iprog.Created[0].Pkg)

			// Build SSA code for bodies of all functions in the whole program.
			prog.Build()

			entryFunc := ssaPkg.Func("main")
			emptyLS := domain.NewEmptyLockSet()
			lsRet, guardedAccessRet := GetFunctionSummary(entryFunc, emptyLS, utils.GetUUID())

			testresult := testutils.TestResult{Lockset: lsRet, GuardedAccess: guardedAccessRet}
			dump, err := json.Marshal(testresult)
			require.NoError(t, err)
			if shouldUpdate {
				utils.UpdateFile(t, tc.resPath, dump)
			}
			expected, err := utils.ReadFile(tc.resPath)
			require.NoError(t, err)

			require.Equal(t, expected, dump)
			_, _ = lsRet, guardedAccessRet

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

			//var labels []string
			for v, l := range result.Queries {
				for _, label := range l.PointsTo().Labels() {
					guardedAccess := valuesQueriesToGuardAccess[v.Pos()]
					allocPos := label.Value()
					positionsToGuardAccesses[allocPos.Pos()] = append(positionsToGuardAccesses[allocPos.Pos()], guardedAccess)
					//label := fmt.Sprintf(" %s with pos:%s may point to %s: %s\n", v, prog.Fset.Position(v.Pos()), prog.Fset.Position(label.Pos()), label)
					//labels = append(labels, label)
				}
			}
			//sort.Strings(labels)
			//for _, label := range labels {
			//	fmt.Println(label)
			//}

			for _, guardedAccesses := range positionsToGuardAccesses {
				for _, guardedAccessesA := range guardedAccesses {
					for _, guardedAccessesB := range guardedAccesses {
						if !guardedAccessesA.Intersects(guardedAccessesB) {
							valueA := guardedAccessesA.Value
							valueB := guardedAccessesB.Value
							label := fmt.Sprintf(" %s with pos:%s has race condition with %s pos:%s \n", valueA, prog.Fset.Position(valueA.Pos()), valueB, prog.Fset.Position(valueB.Pos()))
							print(label)
						}
					}
				}
			}
		})
	}
}
