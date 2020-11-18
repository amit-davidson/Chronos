package ssaUtils

import (
	"github.com/amit-davidson/Chronos/domain"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/ssa"
	"testing"
)

func Test_calculateFunctionStatePathSensitive(t *testing.T) {
	f := LoadMain(t, "./testdata/FunctionsPathSensitive/DeferredLockAndUnlockIfBranch/prog1.go")
	ctx := domain.NewEmptyContext()
	_ = PreProcess
	state := HandleFunction(ctx, f)
	assert.Len(t, state.Lockset.Locks, 0)
	assert.Len(t, state.Lockset.Unlocks, 1)

	block := FindBlockWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
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
	assert.Len(t, block.Lockset.Locks, 0)

	block = FindBlockWithFail(t, state.GuardedAccesses, func(ga *domain.GuardedAccess) bool {
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
	assert.Len(t, block.Lockset.Locks, 0)
}
