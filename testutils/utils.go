package testutils

import (
	"github.com/amit-davidson/Chronos/ssaPureUtils"
	"github.com/amit-davidson/Chronos/ssaUtils"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/ssa"
	"testing"
)

func loadMain(t *testing.T, filePath string) *ssa.Function {
	ssaProg, ssaPkg, err := ssaPureUtils.LoadPackage(filePath)
	require.NoError(t, err)
	f := ssaPkg.Func("main")
	err = ssaUtils.InitPreProcess(ssaProg, ssaPkg, "", f)
	require.NoError(t, err)
	return f
}
