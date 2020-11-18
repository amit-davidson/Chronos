package ssaUtils

import (
	"github.com/amit-davidson/Chronos/domain"
	"github.com/stretchr/testify/assert"
	"go/constant"
	"golang.org/x/tools/go/ssa"
	"testing"
)

func Test_HandleFunction_DeferredLockAndUnlockIfBranch(t *testing.T) {
	f := LoadMain(t, "./testdata/FunctionsPathSensitive/Defer/DeferredLockAndUnlockIfBranch/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 0)
	assert.Len(t, state.Lockset.Unlocks, 1)

	foundGA := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		// find relevant ga
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Const)
		if !ok {
			return false
		}
		if val.Int64() != 5 {
			return false
		}
		return true
	})
	assert.Len(t, foundGA.Lockset.Locks, 0)

	foundGA = FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Const)
		if !ok {
			return false
		}
		if val.Int64() != 6 {
			return false
		}
		return true
	})
	assert.Len(t, foundGA.Lockset.Locks, 0)
}

func Test_HandleFunction_NestedDeferWithLockAndUnlock(t *testing.T) {
	f := LoadMain(t, "./testdata/FunctionsPathSensitive/Defer/NestedDeferWithLockAndUnlock/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 1)
	assert.Len(t, state.Lockset.Unlocks, 0)

	foundGA := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Const)
		if !ok {
			return false
		}
		if val.Int64() != 6 {
			return false
		}
		return true
	})
	assert.Len(t, foundGA.Lockset.Locks, 1)
}

func Test_HandleFunction_NestedDeferWithLockAndUnlockAndGoroutine(t *testing.T) {
	t.Skip("A bug. 7 should contain a lock")
	f := LoadMain(t, "./testdata/FunctionsPathSensitive/Defer/NestedDeferWithLockAndUnlockAndGoroutine/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 0)
	assert.Len(t, state.Lockset.Unlocks, 1)

	foundGA := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Const)
		if !ok {
			return false
		}
		if val.Int64() != 6 {
			return false
		}
		return true
	})
	assert.Len(t, foundGA.Lockset.Locks, 1)

	foundGA = FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Const)
		if !ok {
			return false
		}
		if val.Int64() != 7 {
			return false
		}
		return true
	})
	assert.Len(t, foundGA.Lockset.Locks, 1)
}

func Test_HandleFunction_ForLoop(t *testing.T) {
	t.Skip("A bug. A for loop is assumed to be always executed so state.Lockset.Locks supposed to contain locks")
	f := LoadMain(t, "./testdata/FunctionsPathSensitive/ForLoops/ForLoop/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 1)

	foundGA := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Const)
		if !ok {
			return false
		}
		if constant.StringVal(val.Value) != "b" {
			return false
		}
		return true
	})
	assert.Len(t, foundGA.Lockset.Locks, 0)
	assert.Len(t, foundGA.Lockset.Unlocks, 0)

	foundGA = FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Const)
		if !ok {
			return false
		}
		if constant.StringVal(val.Value) != "c" {
			return false
		}
		return true
	})
	assert.Len(t, foundGA.Lockset.Locks, 1)
	assert.Len(t, foundGA.Lockset.Unlocks, 0)
}

func Test_HandleFunction_NestedForLoopWithRace(t *testing.T) {
	f := LoadMain(t, "./testdata/FunctionsPathSensitive/ForLoops/NestedForLoopWithRace/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 0)
	assert.Len(t, state.Lockset.Unlocks, 0)

	ga := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Const)
		if !ok {
			return false
		}
		if constant.StringVal(val.Value) != "a" {
			return false
		}
		return true
	})
	clockA := ga.State.Clock
	assert.Len(t, ga.Lockset.Locks, 0)
	assert.Len(t, ga.Lockset.Unlocks, 0)

	ga = FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Const)
		if !ok {
			return false
		}
		if constant.StringVal(val.Value) != "b" {
			return false
		}
		return true
	})
	clockB := ga.State.Clock
	assert.Len(t, ga.Lockset.Locks, 0)
	assert.Len(t, ga.Lockset.Unlocks, 0)
	assert.Greater(t, clockB[1], clockA[1])
	assert.Greater(t, clockA[2], clockB[2])
}

func Test_HandleFunction_WhileLoop(t *testing.T) {
	f := LoadMain(t, "./testdata/FunctionsPathSensitive/ForLoops/WhileLoop/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 0)
	assert.Len(t, state.Lockset.Unlocks, 0)

	ga := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		_, ok := ga.Value.(*ssa.Alloc)
		if !ok {
			return false
		}
		pName := ga.Value.Parent().Name()
		if pName == "main" {
			return false
		}
		return true
	})
	clockA := ga.State.Clock
	assert.Len(t, ga.Lockset.Locks, 0)
	assert.Len(t, ga.Lockset.Unlocks, 0)

	ga = FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.FreeVar)
		if !ok {
			return false
		}
		if val.Name() != "x" {
			return false
		}
		return true
	})
	clockB := ga.State.Clock
	assert.Len(t, ga.Lockset.Locks, 0)
	assert.Len(t, ga.Lockset.Unlocks, 0)
	assert.Greater(t, clockB[1], clockA[1])
	assert.Greater(t, clockA[2], clockB[2])
}

func Test_HandleFunction_WhileLoopWithoutHeader(t *testing.T) {
	t.Skip("for {}")
	f := LoadMain(t, "./testdata/FunctionsPathSensitive/ForLoops/WhileLoopWithoutHeader/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 0)
	assert.Len(t, state.Lockset.Unlocks, 0)

	ga := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		_, ok := ga.Value.(*ssa.Alloc)
		if !ok {
			return false
		}
		pName := ga.Value.Parent().Name()
		if pName == "main" {
			return false
		}
		return true
	})
	clockA := ga.State.Clock
	assert.Len(t, ga.Lockset.Locks, 0)
	assert.Len(t, ga.Lockset.Unlocks, 0)

	ga = FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.FreeVar)
		if !ok {
			return false
		}
		if val.Name() != "x" {
			return false
		}
		return true
	})
	clockB := ga.State.Clock
	assert.Len(t, ga.Lockset.Locks, 0)
	assert.Len(t, ga.Lockset.Unlocks, 0)
	assert.Greater(t, clockB[1], clockA[1])
	assert.Greater(t, clockA[2], clockB[2])
}
