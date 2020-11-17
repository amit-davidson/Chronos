package ssaUtils

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/ssa"
	"testing"
)

func LoadMain(t *testing.T, filePath string) *ssa.Function {
	ssaProg, ssaPkg, err := LoadPackage(filePath)
	require.NoError(t, err)
	f := ssaPkg.Func("main")
	err = InitPreProcess(ssaProg, ssaPkg, "", f)
	require.NoError(t, err)
	return f
}
