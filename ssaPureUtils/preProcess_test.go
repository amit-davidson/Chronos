package ssaPureUtils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func setup(t *testing.T, filePath string) bool {
	ssaProg, ssaPkg, err := LoadPackage(filePath)
	require.NoError(t, err)
	f := ssaPkg.Func("main")
	PreProcess := InitPreProcess(f)
	err = SetGlobals(ssaProg, ssaPkg, PreProcess, "")
	require.NoError(t, err)
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	return isContainingLock
}

func TestIsFunctionContainingLocks_Empty(t *testing.T) {
	isContainingLock := setup(t, "./testdata/functionWithoutLocks/functionWithoutLocks.go")
	require.False(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithLock(t *testing.T) {
	isContainingLock := setup(t, "./testdata/functionWithLock/functionWithLock.go")
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithUnlock(t *testing.T) {
	isContainingLock := setup(t, "./testdata/functionWithUnlock/functionWithUnlock.go")
	require.False(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithEvenLockAndUnlock(t *testing.T) {
	isContainingLock := setup(t, "./testdata/functionWithEvenLockAndUnlock/functionWithEvenLockAndUnlock.go")
	require.False(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithEvenLockAndUnlockAndLockComesLater(t *testing.T) {
	// Even though the number of locks and unlocks is equal, this still marked as containing locks since in the data
	// flow analysis, the unlock precede the lock.
	isContainingLock := setup(t, "./testdata/EvenLockAndUnlockAndLockComesLater/EvenLockAndUnlockAndLockComesLater.go")
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithEvenLockAndUnlockAndLockComesLaterInTheFlowButAppearEarlierInCode(t *testing.T) {
	// In general this should be marked as true since the mutex is marked as locked when the function exits. It's
	// assumed that if the number of locks equals to the number of unlocks, then the lock step will precede the unlock step.
	// Even if the number of cases is equal and the order is opposite, the data flow analysis still might see it
	// due to the order they appear in the code. See TestIsFunctionContainingLocks_WithEvenLockAndUnlockAndLockComesLater for example.
	// Cases as in this may result in incorrect/missing reports of race conditions but they are very rare.
	isContainingLock := setup(t, "./testdata/TestIsFunctionContainingLocks_WithEvenLockAndUnlockAndLockComesLaterInTheFlowButAppearEarlierInCode/TestIsFunctionContainingLocks_WithEvenLockAndUnlockAndLockComesLaterInTheFlowButAppearEarlierInCode.go")
	require.False(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithEvenLockAndDeferUnlock(t *testing.T) {
	isContainingLock := setup(t, "./testdata/functionWithEvenLockAndDeferUnlock/functionWithEvenLockAndDeferUnlock.go")
	require.False(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithUnevenLockAndUnlock(t *testing.T) {
	// In general this should be marked as false since the mutex will always be marked as unlocked when the function exits.
	// It's usually appear when a lock on a mutex might appear under multiple branches but a single unlock to ensure the mutex is released.
	// In the test, a full analysis of the control flow graph is not needed but the only harm is performance.
	isContainingLock := setup(t, "./testdata/functionWithUnevenLockAndUnlock/functionWithUnevenLockAndUnlock.go")
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithUnevenLockAndUnlockAndLockComesLater(t *testing.T) {
	isContainingLock := setup(t, "./testdata/UnevenLockAndUnlockAndLockComesLater/UnevenLockAndUnlockAndLockComesLater.go")
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_NestedFunctionWithLock(t *testing.T) {
	isContainingLock := setup(t, "./testdata/NestedFunctionWithLock/NestedFunctionWithLock.go")
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_NestedFunctionWithLockButThenUnlock(t *testing.T) {
	isContainingLock := setup(t, "./testdata/NestedFunctionWithLockButThenUnlock/NestedFunctionWithLockButThenUnlock.go")
	require.False(t, isContainingLock)
}