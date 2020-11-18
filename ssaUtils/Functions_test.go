package ssaUtils

import (
	"github.com/amit-davidson/Chronos/domain"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/ssa"
	"testing"
)

func Test_HandleFunction_DeferredLockAndUnlockIfBranch(t *testing.T) {
	f := LoadMain(t, "./testdata/Functions/Defer/DeferredLockAndUnlockIfBranch/prog1.go")
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
	f := LoadMain(t, "./testdata/Functions/Defer/NestedDeferWithLockAndUnlock/prog1.go")
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
	f := LoadMain(t, "./testdata/Functions/Defer/NestedDeferWithLockAndUnlockAndGoroutine/prog1.go")
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
	f := LoadMain(t, "./testdata/Functions/ForLoops/ForLoop/prog1.go")
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
	f := LoadMain(t, "./testdata/Functions/ForLoops/NestedForLoopWithRace/prog1.go")
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
	f := LoadMain(t, "./testdata/Functions/ForLoops/WhileLoop/prog1.go")
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
	t.Skip("for {}")
	f := LoadMain(t, "./testdata/Functions/ForLoops/WhileLoopWithoutHeader/prog1.go")
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

func Test_HandleFunction_DataRaceIceCreamMaker(t *testing.T) {
	f := LoadMain(t, "./testdata/Functions/Interfaces/DataRaceIceCreamMaker/prog1.go")
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
	f := LoadMain(t, "./testdata/Functions/Interfaces/InterfaceWithLock/prog1.go")
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
	f := LoadMain(t, "./testdata/Functions/Interfaces/NestedInterface/prog1.go")
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
	f := LoadMain(t, "./testdata/Functions/LocksAndUnlocks/Lock/prog1.go")
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
	f := LoadMain(t, "./testdata/Functions/LocksAndUnlocks/LockAndUnlock/prog1.go")
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
	f := LoadMain(t, "./testdata/Functions/LocksAndUnlocks/LockAndUnlockIfBranch/prog1.go")
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
	f := LoadMain(t, "./testdata/Functions/LocksAndUnlocks/LockInBothBranches/prog1.go")
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
	f := LoadMain(t, "./testdata/Functions/LocksAndUnlocks/LockInsideGoroutine/prog1.go")
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
	f := LoadMain(t, "./testdata/Functions/LocksAndUnlocks/MultipleLocksNoRace/prog1.go")
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
	assert.True(t, GA1.State.MayConcurrent(GA2.State))
	assert.True(t, GA1.Intersects(GA2)) // Locks intersect
}

func Test_HandleFunction_NestedConditionWithLockInAllBranches(t *testing.T) {
	f := LoadMain(t, "./testdata/Functions/LocksAndUnlocks/NestedConditionWithLockInAllBranches/prog1.go")
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
	f := LoadMain(t, "./testdata/Functions/LocksAndUnlocks/NestedLockInStruct/prog1.go")
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
