package ssaUtils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIsFunctionContainingLocks_Empty(t *testing.T) {
	f := LoadMain(t, "./testdata/Preprocess/functionWithoutLocks/functionWithoutLocks.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.False(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithLock(t *testing.T) {
	f := LoadMain(t, "./testdata/Preprocess/functionWithLock/functionWithLock.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithUnlock(t *testing.T) {
	f := LoadMain(t, "./testdata/Preprocess/functionWithUnlock/functionWithUnlock.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.False(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithLockAndUnlockInDifferentBranches(t *testing.T) {
	// The approximation here is not correct. Full analysis of the cfg is not needed here since not all paths contain
	// locks, but it misses and full search will be performed. That's deliberate and a trade-off.
	f := LoadMain(t, "./testdata/Preprocess/LockAndUnlockInDifferentBranches/LockAndUnlockInDifferentBranches.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithEvenLockAndUnlock(t *testing.T) {
	f := LoadMain(t, "./testdata/Preprocess/functionWithEvenLockAndUnlock/functionWithEvenLockAndUnlock.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.False(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithEvenLockAndUnlockAndLockComesLater(t *testing.T) {
	f := LoadMain(t, "./testdata/Preprocess/EvenLockAndUnlockAndLockComesLater/LockInEmbeddedStruct.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithEvenLockAndDeferUnlock(t *testing.T) {
	f := LoadMain(t, "./testdata/Preprocess/functionWithEvenLockAndDeferUnlock/functionWithEvenLockAndDeferUnlock.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.False(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithUnevenLockAndUnlock(t *testing.T) {
	f := LoadMain(t, "./testdata/Preprocess/functionWithUnevenLockAndUnlock/functionWithUnevenLockAndUnlock.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.False(t, isContainingLock)
}

func TestIsFunctionContainingLocks_WithUnevenLockAndUnlockAndLockComesLater(t *testing.T) {
	f := LoadMain(t, "./testdata/Preprocess/UnevenLockAndUnlockAndLockComesLater/WaitGroup.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.False(t, isContainingLock)
}

func TestIsFunctionContainingLocks_NestedFunctionWithLock(t *testing.T) {
	f := LoadMain(t, "./testdata/Preprocess/NestedFunctionWithLock/NestedFunctionWithLock.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_RecursionNestedFunctionWithLock(t *testing.T) {
	f := LoadMain(t, "./testdata/Preprocess/RecursionNestedFunctionWithLock/RecursionNestedFunctionWithLock.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_NestedFunctionWithLockButThenUnlock(t *testing.T) {
	f := LoadMain(t, "./testdata/Preprocess/NestedFunctionWithLockButThenUnlock/NestedFunctionWithLockButThenUnlock.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.False(t, isContainingLock)
}

func TestIsFunctionContainingLocks_LockInAStruct(t *testing.T) {
	f := LoadMain(t, "./testdata/Preprocess/LockInAStruct/LockInEmbeddedStruct.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.True(t, isContainingLock)
}

func TestIsFunctionContainingLocks_EmbeddedStruct(t *testing.T) {
	f := LoadMain(t, "./testdata/Preprocess/LockInEmbeddedStruct/LockInEmbeddedStruct.go")
	isContainingLock := PreProcess.FunctionWithLocks[f.Signature]
	require.True(t, isContainingLock)
}