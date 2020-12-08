package ssaUtils

import (
	"go/constant"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/ssa"
)

func FindGA(GuardedAccesses []*domain.GuardedAccess, validationFunc func(value *domain.GuardedAccess) bool) *domain.GuardedAccess {
	wasFound := false
	for _, ga := range GuardedAccesses {
		wasFound = validationFunc(ga)
		if wasFound == true {
			return ga
		}
	}
	return nil
}

func FindMultipleGA(GuardedAccesses []*domain.GuardedAccess, validationFunc func(value *domain.GuardedAccess) bool) []*domain.GuardedAccess {
	foundGAs := make([]*domain.GuardedAccess, 0)
	for _, ga := range GuardedAccesses {
		wasFound := validationFunc(ga)
		if wasFound == true {
			foundGAs = append(foundGAs, ga)
		}
	}
	return foundGAs
}

func GetConstString(v *ssa.Const) string {
	return constant.StringVal(v.Value)
}

func GetGlobalString(v *ssa.Global) string {
	return v.Name()
}

func IsGARead(ga *domain.GuardedAccess) bool {
	return ga.OpKind == domain.GuardAccessRead
}

func IsGAWrite(ga *domain.GuardedAccess) bool {
	return ga.OpKind == domain.GuardAccessWrite
}

func FindGAWithFail(t *testing.T, GuardedAccesses []*domain.GuardedAccess, validationFunc func(value *domain.GuardedAccess) bool) *domain.GuardedAccess {
	res := FindGA(GuardedAccesses, validationFunc)
	require.NotNil(t, res)
	return res
}

func FindMultipleGAWithFail(t *testing.T, GuardedAccesses []*domain.GuardedAccess, validationFunc func(value *domain.GuardedAccess) bool, expectedAmount int) []*domain.GuardedAccess {
	res := FindMultipleGA(GuardedAccesses, validationFunc)
	require.Equal(t, expectedAmount, len(res))
	return res
}

func LoadMain(t *testing.T, filePath string) (*ssa.Function, *ssa.Package) {
	domain.GoroutineCounter = utils.NewCounter()
	domain.GuardedAccessCounter = utils.NewCounter()
	domain.PosIDCounter = utils.NewCounter()

	_, ex, _, ok := runtime.Caller(0)
	require.True(t, ok)
	modulePath := filepath.Dir(filepath.Dir(ex))

	ssaProg, ssaPkg, err := LoadPackage(filePath, modulePath)
	require.NoError(t, err)
	f := ssaPkg.Func("main")
	err = InitPreProcess(ssaProg, modulePath)
	require.NoError(t, err)
	return f, ssaPkg
}

func EqualDifferentOrder(a, b []*domain.GuardedAccess) bool {
	if len(a) != len(b) {
		return false
	}
	diff := make(map[int]int, len(a))
	for _, x := range a {
		diff[x.ID]++
	}
	for _, y := range b {
		if _, ok := diff[y.ID]; !ok {
			return false
		}
		diff[y.ID] -= 1
		if diff[y.ID] == 0 {
			delete(diff, y.ID)
		}
	}
	return len(diff) == 0
}
