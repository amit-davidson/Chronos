package ssaUtils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsFunctionContainingLocks_Empty(t *testing.T) {
	f := LoadMain(t, "./testdata/functionWithoutLocks/functionWithoutLocks.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.False(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithLock(t *testing.T) {
	f := LoadMain(t, "./testdata/functionWithLock/functionWithLock.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithUnlock(t *testing.T) {
	f := LoadMain(t, "./testdata/functionWithUnlock/functionWithUnlock.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.False(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithEvenLockAndUnlock(t *testing.T) {
	f := LoadMain(t, "./testdata/functionWithEvenLockAndUnlock/functionWithEvenLockAndUnlock.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.False(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithEvenLockAndUnlockAndLockComesLater(t *testing.T) {
	// Even though the number of locks and unlocks is equal, this still marked as containing locks since in the data
	// flow analysis, the unlock precedes the lock.
	f := LoadMain(t, "./testdata/EvenLockAndUnlockAndLockComesLater/EvenLockAndUnlockAndLockComesLater.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithEvenLockAndDeferUnlock(t *testing.T) {
	f := LoadMain(t, "./testdata/functionWithEvenLockAndDeferUnlock/functionWithEvenLockAndDeferUnlock.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.False(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithUnevenLockAndUnlock(t *testing.T) {
	// In general this should be marked as false since the mutex will always be marked as unlocked when the function exits.
	// It's usually appear when a lock on a mutex might appear under multiple branches but a single unlock to ensure the mutex is released.
	// In the test, a full analysis of the control flow graph is not needed but the only harm is performance.
	f := LoadMain(t, "./testdata/functionWithUnevenLockAndUnlock/functionWithUnevenLockAndUnlock.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithUnevenLockAndUnlockAndLockComesLater(t *testing.T) {
	f := LoadMain(t, "./testdata/UnevenLockAndUnlockAndLockComesLater/UnevenLockAndUnlockAndLockComesLater.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_NestedFunctionWithLock(t *testing.T) {
	f := LoadMain(t, "./testdata/NestedFunctionWithLock/NestedFunctionWithLock.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_RecursionNestedFunctionWithLock(t *testing.T) {
	f := LoadMain(t, "./testdata/RecursionNestedFunctionWithLock/RecursionNestedFunctionWithLock.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_NestedFunctionWithLockButThenUnlock(t *testing.T) {
	f := LoadMain(t, "./testdata/NestedFunctionWithLockButThenUnlock/NestedFunctionWithLockButThenUnlock.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.False(t, isContainingLock)
}
