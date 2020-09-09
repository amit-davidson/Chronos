package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
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
			emptyLS := newEmptyLockSet()
			lsRet, guardedAccessRet := GetFunctionSummary(entryFunc, emptyLS)
			dumpLs, err := lsRet.MarshalJSON()
			require.NoError(t, err)
			for _, guardedAccess := range guardedAccessRet {
				dumpGuardedAccess, err := guardedAccess.MarshalJSON()
				require.NoError(t, err)
				dumpLs = append(dumpLs, []byte{'\n'}...)
				dumpLs = append(dumpLs, dumpGuardedAccess...)
			}
			UpdateFile(t, tc.resPath, dumpLs, shouldUpdate)
			expected, err := ReadFile(tc.resPath)
			require.NoError(t, err)

			require.Equal(t, expected, dumpLs)
			_, _ = lsRet, guardedAccessRet

			config := &pointer.Config{
				Mains: []*ssa.Package{ssaPkg},
			}

			for _, guardedAccess := range guardedAccessRet {
				if pointer.CanPoint(guardedAccess.value.Type()) {
					config.AddQuery(guardedAccess.value)
				}
			}
			result, err := pointer.Analyze(config)
			if err != nil {
				panic(err) // internal error in pointer analysis
			}
			var labels []string
			for v, l := range result.Queries {
				for _, t := range l.PointsTo().Labels() {
					//name := config.Queries[l]
					label := fmt.Sprintf(" %s with pos:%s may point to %s: %s\n", v, prog.Fset.Position(v.Pos()), prog.Fset.Position(t.Pos()), t)
					labels = append(labels, label)
				}
			}
			//sort.Strings(labels)
			for _, label := range labels {
				fmt.Println(label)
			}

		})
	}
}
