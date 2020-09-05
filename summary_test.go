package main

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/loader"
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
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var conf loader.Config
			file, err := conf.ParseFile(tc.testPath, nil)
			conf.CreateFromFiles("main", file)

			iprog, err := conf.Load()
			require.NoError(t, err) // Type error in some package

			prog := ssautil.CreateProgram(iprog, 0)
			mainPkg := prog.Package(iprog.Created[0].Pkg)

			prog.Build()

			funcInterface5 := mainPkg.Func("fn1")
			lsRet, guardedAccessRet := GetFunctionSummary(funcInterface5)
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
		})
	}
}
