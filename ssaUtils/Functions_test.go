package ssaUtils

import (
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/pointerAnalysis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/ssa"
	"testing"
)

func Test_HandleFunction_DeferredLockAndUnlockIfBranch(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/Defer/DeferredLockAndUnlockIfBranch/prog1.go")
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
	f, _ := LoadMain(t, "./testdata/Functions/Defer/NestedDeferWithLockAndUnlock/prog1.go")
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
	//t.Skip("A bug. 7 should contain a lock")
	f, _ := LoadMain(t, "./testdata/Functions/Defer/NestedDeferWithLockAndUnlockAndGoroutine/prog1.go")
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

func Test_HandleFunction_ForLoopLockInsideLoop(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/ForLoops/ForLoopLockInsideLoop/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 0)

	foundGA := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Const)
		if !ok {
			return false
		}
		if GetConstString(val) != "b" {
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
		if GetConstString(val) != "c" {
			return false
		}
		return true
	})
	assert.Len(t, foundGA.Lockset.Locks, 1)
	assert.Len(t, foundGA.Lockset.Unlocks, 0)
}

func Test_HandleFunction_ForLoopLockOutsideLoop(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/ForLoops/ForLoopLockOutsideLoop/prog1.go")
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
		if GetConstString(val) != "b" {
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
		if GetConstString(val) != "c" {
			return false
		}
		return true
	})
	assert.Len(t, foundGA.Lockset.Locks, 1)
	assert.Len(t, foundGA.Lockset.Unlocks, 0)
}

func Test_HandleFunction_NestedForLoopWithRace(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/ForLoops/NestedForLoopWithRace/prog1.go")
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
		if GetConstString(val) != "a" {
			return false
		}
		return true
	})
	stateA := ga.State
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
		if GetConstString(val) != "b" {
			return false
		}
		return true
	})
	stateB := ga.State
	assert.Len(t, ga.Lockset.Locks, 0)
	assert.Len(t, ga.Lockset.Unlocks, 0)
	assert.True(t, stateA.MayConcurrent(stateB))

}

func Test_HandleFunction_WhileLoop(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/ForLoops/WhileLoop/prog1.go")
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
	stateA := ga.State
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
	stateB := ga.State
	assert.Len(t, ga.Lockset.Locks, 0)
	assert.Len(t, ga.Lockset.Unlocks, 0)
	assert.True(t, stateA.MayConcurrent(stateB))

}

func Test_HandleFunction_WhileLoopWithoutHeader(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/ForLoops/WhileLoopWithoutHeader/prog1.go")
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
	stateA := ga.State
	assert.Len(t, ga.Lockset.Locks, 1)
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
	stateB := ga.State
	assert.Len(t, ga.Lockset.Locks, 0)
	assert.Len(t, ga.Lockset.Unlocks, 0)
	assert.True(t, stateA.MayConcurrent(stateB))
}

func Test_HandleFunction_DataRaceIceCreamMaker(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/Interfaces/DataRaceIceCreamMaker/prog1.go")
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
		if GetConstString(val) != "Ben" {
			return false
		}
		return true
	})
	stateA := ga.State
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
		if GetConstString(val) != "1" {
			return false
		}
		return true
	})
	stateB := ga.State
	assert.Len(t, ga.Lockset.Locks, 0)
	assert.Len(t, ga.Lockset.Unlocks, 0)
	assert.True(t, stateA.MayConcurrent(stateB))
}

func Test_HandleFunction_InterfaceWithLock(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/Interfaces/InterfaceWithLock/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 1)
	assert.Len(t, state.Lockset.Unlocks, 0)

	ga := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Const)
		if !ok {
			return false
		}
		if GetConstString(val) != "Ben" {
			return false
		}
		return true
	})
	assert.Len(t, ga.Lockset.Locks, 1)
	assert.Len(t, ga.Lockset.Unlocks, 0)

	ga = FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Const)
		if !ok {
			return false
		}
		if GetConstString(val) != "Jerry" {
			return false
		}
		return true
	})
	assert.Len(t, ga.Lockset.Locks, 1)
	assert.Len(t, ga.Lockset.Unlocks, 0)

	ga = FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Const)
		if !ok {
			return false
		}
		if GetConstString(val) != "1" {
			return false
		}
		return true
	})
	assert.Len(t, ga.Lockset.Locks, 1)
	assert.Len(t, ga.Lockset.Unlocks, 0)
}

func Test_HandleFunction_NestedInterface(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/Interfaces/NestedInterface/prog1.go")
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
		if GetConstString(val) != "Jerry" {
			return false
		}
		return true
	})
	assert.Len(t, ga.Lockset.Locks, 0)
	assert.Len(t, ga.Lockset.Unlocks, 0)
}

func Test_HandleFunction_Lock(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/LocksAndUnlocks/Lock/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 1)
	assert.Len(t, state.Lockset.Unlocks, 0)

	ga := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
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
	assert.Len(t, ga.Lockset.Locks, 1)
	assert.Len(t, ga.Lockset.Unlocks, 0)
}

func Test_HandleFunction_LockAndUnlock(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/LocksAndUnlocks/LockAndUnlock/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 0)
	assert.Len(t, state.Lockset.Unlocks, 1)

	ga := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
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
	assert.Len(t, ga.Lockset.Locks, 0)
	assert.Len(t, ga.Lockset.Unlocks, 1)
}

func Test_HandleFunction_LockAndUnlockIfBranch(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/LocksAndUnlocks/LockAndUnlockIfBranch/prog1.go")
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

func Test_HandleFunction_LockInBothBranches(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/LocksAndUnlocks/LockInBothBranches/prog1.go")
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
		if val.Int64() != 5 {
			return false
		}
		return true
	})
	assert.Len(t, foundGA.Lockset.Locks, 1)
}

func Test_HandleFunction_LockInsideGoroutine(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/LocksAndUnlocks/LockInsideGoroutine/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 0)
	assert.Len(t, state.Lockset.Unlocks, 0)

	foundGA := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		_, ok := ga.Value.(*ssa.MakeInterface)
		if !ok {
			return false
		}
		pName := ga.Value.Parent().Name()
		if pName == "main" {
			return false
		}
		return true
	})
	assert.Len(t, foundGA.Lockset.Locks, 1)

	foundGA = FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		_, ok := ga.Value.(*ssa.MakeInterface)
		if !ok {
			return false
		}
		pName := ga.Value.Parent().Name()
		if pName != "main" {
			return false
		}
		return true
	})
	assert.Len(t, foundGA.Lockset.Locks, 0)
}

func Test_HandleFunction_MultipleLocksNoRace(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/LocksAndUnlocks/MultipleLocksNoRace/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 0)
	assert.Len(t, state.Lockset.Unlocks, 0)

	GA1 := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Const)
		if !ok {
			return false
		}
		if val.Int64() != 1 {
			return false
		}
		return true
	})
	assert.Len(t, GA1.Lockset.Locks, 1)

	GA2 := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Const)
		if !ok {
			return false
		}
		if val.Int64() != 2 {
			return false
		}
		return true
	})

	assert.Len(t, GA2.Lockset.Locks, 1)
	assert.False(t, GA1.IsConflicting(GA2))
}

func Test_HandleFunction_NestedConditionWithLockInAllBranches(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/LocksAndUnlocks/NestedConditionWithLockInAllBranches/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 1)
	assert.Len(t, state.Lockset.Unlocks, 0)

	GA1 := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		_, ok := ga.Value.(*ssa.MakeInterface)
		if !ok {
			return false
		}
		return true
	})
	assert.Len(t, GA1.Lockset.Locks, 1)
}

func Test_HandleFunction_NestedLockInStruct(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/LocksAndUnlocks/NestedLockInStruct/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 0)
	assert.Len(t, state.Lockset.Unlocks, 1)

	GA1 := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
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
	assert.Len(t, GA1.Lockset.Locks, 0)
	assert.Len(t, GA1.Lockset.Unlocks, 1)
}

func Test_HandleFunction_DataRaceGoto(t *testing.T) {
	f, pkg := LoadMain(t, "./testdata/Functions/General/DataRaceGoto/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	conflictingAccesses, err := pointerAnalysis.Analysis(pkg, state.GuardedAccesses)
	require.NoError(t, err)
	gas := FindMultipleGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGAWrite(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Global)
		if !ok {
			return false
		}
		if GetGlobalString(val) != "a" {
			return false
		}
		return true
	}, 2)

	found := false
	for _, ca := range conflictingAccesses {
		if EqualDifferentOrder(gas, ca) {
			found = true
		}
	}
	assert.True(t, found)

	ga := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGAWrite(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Global)
		if !ok {
			return false
		}
		if GetGlobalString(val) != "a" {
			return false
		}
		return true
	})

	found = false
	for _, ca := range conflictingAccesses {
		if EqualDifferentOrder([]*domain.GuardedAccess{gas[1], ga}, ca) {
			found = true
		}
	}
	assert.True(t, found)
}

func Test_HandleFunction_DataRaceMap(t *testing.T) {
	f, pkg := LoadMain(t, "./testdata/Functions/General/DataRaceMap/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	conflictingAccesses, err := pointerAnalysis.Analysis(pkg, state.GuardedAccesses)
	require.NoError(t, err)
	gaA := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGAWrite(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.UnOp)
		if !ok {
			return false
		}
		X, ok := val.X.(*ssa.FreeVar)
		if !ok {
			return false
		}
		if X.Name() != "m" {
			return false
		}
		return true
	})

	gaB := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGAWrite(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.UnOp)
		if !ok {
			return false
		}
		X, ok := val.X.(*ssa.Alloc)
		if !ok {
			return false
		}
		if X.Comment != "m" {
			return false
		}
		return true
	})

	found := false
	for _, ca := range conflictingAccesses {
		if EqualDifferentOrder([]*domain.GuardedAccess{gaA, gaB}, ca) {
			found = true
		}
	}
	assert.True(t, found)
}

func Test_HandleFunction_DataRaceNestedSameFunction(t *testing.T) {
	f, pkg := LoadMain(t, "./testdata/Functions/General/DataRaceNestedSameFunction/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	conflictingAccesses, err := pointerAnalysis.Analysis(pkg, state.GuardedAccesses)
	require.NoError(t, err)
	gas := FindMultipleGA(state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		global, ok := ga.Value.(*ssa.Global)
		if !ok {
			return false
		}
		if global.Name() != "count" {
			return false
		}
		return true
	})

	assert.Subset(t, gas, conflictingAccesses[0])
}

func Test_HandleFunction_DataRaceProperty(t *testing.T) {
	f, pkg := LoadMain(t, "./testdata/Functions/General/DataRaceProperty/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	conflictingAccesses, err := pointerAnalysis.Analysis(pkg, state.GuardedAccesses)
	require.NoError(t, err)
	gaA := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGAWrite(ga) {
			return false
		}
		field, ok := ga.Value.(*ssa.FieldAddr)
		if !ok {
			return false
		}
		if field.X.Name() != "w" {
			return false
		}
		return true
	})
	assert.Len(t, gaA.State.Clock, 2)

	gaB := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		field, ok := ga.Value.(*ssa.FieldAddr)
		if !ok {
			return false
		}
		_, ok = field.X.(*ssa.UnOp)
		if !ok {
			return false
		}
		return true
	})
	assert.Len(t, gaB.State.Clock, 3)

	found := false
	for _, ca := range conflictingAccesses {
		if EqualDifferentOrder([]*domain.GuardedAccess{gaA, gaB}, ca) {
			found = true
		}
	}
	assert.True(t, found)
}

func Test_HandleFunction_DataRaceRecursion(t *testing.T) {
	t.Skip("gas should be of length 3")
	f, pkg := LoadMain(t, "./testdata/Functions/General/DataRaceRecursion/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	conflictingAccesses, err := pointerAnalysis.Analysis(pkg, state.GuardedAccesses)
	require.NoError(t, err)
	gas := FindMultipleGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGAWrite(ga) {
			return false
		}
		field, ok := ga.Value.(*ssa.Global)
		if !ok {
			return false
		}
		if field.Name() != "a" {
			return false
		}
		return true
	}, 2)

	found := false
	for _, ca := range conflictingAccesses {
		if EqualDifferentOrder(gas, ca) {
			found = true
		}
	}
	assert.True(t, found)
}

func Test_HandleFunction_DataRaceShadowedErr(t *testing.T) {
	f, pkg := LoadMain(t, "./testdata/Functions/General/DataRaceShadowedErr/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	conflictingAccesses, err := pointerAnalysis.Analysis(pkg, state.GuardedAccesses)
	require.NoError(t, err)
	assert.Len(t, conflictingAccesses, 6)
	filteredAccesses := pointerAnalysis.FilterDuplicates(conflictingAccesses)
	assert.Len(t, filteredAccesses, 3)
}

func Test_HandleFunction_DataRaceWithOnlyAlloc(t *testing.T) {
	f, pkg := LoadMain(t, "./testdata/Functions/General/DataRaceWithOnlyAlloc/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	conflictingAccesses, err := pointerAnalysis.Analysis(pkg, state.GuardedAccesses)
	require.NoError(t, err)
	assert.Len(t, conflictingAccesses, 2)
	filteredAccesses := pointerAnalysis.FilterDuplicates(conflictingAccesses)
	assert.Len(t, filteredAccesses, 1)

	gaA := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGAWrite(ga) {
			return false
		}
		_, ok := ga.Value.(*ssa.FreeVar)
		if !ok {
			return false
		}
		return true
	})

	gaB := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		_, ok := ga.Value.(*ssa.Alloc)
		if !ok {
			return false
		}
		return true
	})
	require.Len(t, filteredAccesses, 1)
	assert.True(t, EqualDifferentOrder([]*domain.GuardedAccess{gaA, gaB}, filteredAccesses[0]))
}

func Test_HandleFunction_DataRaceWithSameFunction(t *testing.T) {
	f, pkg := LoadMain(t, "./testdata/Functions/General/DataRaceWithSameFunction/prog1.go")
	ctx := domain.NewEmptyContext()
	entryCallCommon := ssa.CallCommon{Value: f}
	state := HandleCallCommon(ctx, &entryCallCommon, f.Pos())
	conflictingAccesses, err := pointerAnalysis.Analysis(pkg, state.GuardedAccesses)
	require.NoError(t, err)
	filteredAccesses := pointerAnalysis.FilterDuplicates(conflictingAccesses)
	assert.Len(t, filteredAccesses, 1)

	gas := FindMultipleGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		val, ok := ga.Value.(*ssa.Global)
		if !ok {
			return false
		}
		if GetGlobalString(val) != "count" {
			return false
		}
		return true
	}, 4)
	for _, ga := range gas {
		assert.Len(t, ga.State.StackTrace.Iter(), 2)
	}
	require.Len(t, filteredAccesses, 1)
	assert.Subset(t, gas, conflictingAccesses[0])
}

func Test_HandleFunction_NestedFunctions(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/General/NestedFunctions/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 1)
	assert.Len(t, state.Lockset.Unlocks, 2)

	ga := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
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
	assert.Len(t, ga.Lockset.Locks, 2)
	assert.Len(t, ga.Lockset.Unlocks, 0)

	ga = FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
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
	assert.Len(t, ga.Lockset.Locks, 3)
	assert.Len(t, ga.Lockset.Unlocks, 0)
}

func Test_HandleFunction_RecursionWithGoroutine(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/General/RecursionWithGoroutine/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)

	// Should found 2 occurrences since the algorithm should traverse more then once to find conflicting accesses such
	// as recursion with a goroutine
	_ = FindMultipleGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.Const)
		if !ok {
			return false
		}
		if val.Int64() != 0 {
			return false
		}
		return true
	}, 2)
}

func Test_HandleFunction_Simple(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/General/Simple/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)

	gas := FindMultipleGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGAWrite(ga) {
			return false
		}
		_, ok := ga.Value.(*ssa.FieldAddr)
		if !ok {
			return false
		}
		return true
	}, 2)
	assert.True(t, gas[0].State.MayConcurrent(gas[1].State))
}

func Test_HandleFunction_StructMethod(t *testing.T) {
	f, _ := LoadMain(t, "./testdata/Functions/General/StructMethod/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)

	gaA := FindMultipleGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGAWrite(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.FieldAddr)
		if !ok {
			return false
		}
		_, ok = val.X.(*ssa.Alloc)
		if ok {
			return false
		}
		return true
	}, 1)

	gaB := FindGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGARead(ga) {
			return false
		}
		val, ok := ga.Value.(*ssa.FieldAddr)
		if !ok {
			return false
		}
		_, ok = val.X.(*ssa.Parameter)
		if !ok {
			return false
		}
		return true
	})
	assert.True(t, gaA[0].State.MayConcurrent(gaB.State))
}

func Test_HandleFunction_DataRaceInterfaceOverChannel(t *testing.T) {
	f, pkg := LoadMain(t, "./testdata/Functions/PointerAnalysis/DataRaceInterfaceOverChannel/prog1.go")
	ctx := domain.NewEmptyContext()
	state := HandleFunction(ctx, f)
	

	gas := FindMultipleGAWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
		if !IsGAWrite(ga) {
			return false
		}
		_, ok := ga.Value.(*ssa.FieldAddr)
		if !ok {
			return false
		}
		return true
	}, 2)


	conflictingAccesses, err := pointerAnalysis.Analysis(pkg, state.GuardedAccesses)
	require.NoError(t, err)
	filteredAccesses := pointerAnalysis.FilterDuplicates(conflictingAccesses)
	require.Len(t, filteredAccesses, 1)
	found := false
	for _, ca := range conflictingAccesses {
		if EqualDifferentOrder(gas, ca) {
			found = true
		}
	}
	assert.True(t, found)
}
