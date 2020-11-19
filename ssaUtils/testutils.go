package ssaUtils

import (
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/utils"
	"github.com/stretchr/testify/require"
	"go/constant"
	"golang.org/x/tools/go/ssa"
	"testing"
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

func GetConstString(v *ssa.Const) string {
	return constant.StringVal(v.Value)

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

func LoadMain(t *testing.T, filePath string) *ssa.Function {
	domain.GoroutineCounter = utils.NewCounter()
	domain.GuardedAccessCounter = utils.NewCounter()
	domain.PosIDCounter = utils.NewCounter()

	ssaProg, ssaPkg, err := LoadPackage(filePath)
	require.NoError(t, err)
	f := ssaPkg.Func("main")
	err = InitPreProcess(ssaProg, ssaPkg, "", f)
	require.NoError(t, err)
	return f
}