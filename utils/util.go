package utils

import (
	"bufio"
	"github.com/stretchr/testify/require"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ssa"
	"io/ioutil"
	"os"
	"testing"
)

func IsCallToAny(call *ssa.CallCommon, names ...string) bool {
	q := CallName(call, false)
	for _, name := range names {
		if q == name {
			return true
		}
	}
	return false
}

func CallName(call *ssa.CallCommon, short bool) string {
	if call.IsInvoke() {
		return ""
	}
	switch v := call.Value.(type) {
	case *ssa.Function:
		fn, ok := v.Object().(*types.Func)
		if !ok {
			return ""
		}
		if short {
			return fn.Name()
		} else {
			return fn.FullName()
		}
	case *ssa.Builtin:
		return v.Name()
	}
	return ""
}

func FilterDebug(instr []ssa.Instruction) []ssa.Instruction {
	var out []ssa.Instruction
	for _, ins := range instr {
		if _, ok := ins.(*ssa.DebugRef); !ok {
			out = append(out, ins)
		}
	}
	return out
}

func OpenFile(fileName string) (*os.File, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func CreateFile(fileName string) (*os.File, error) {
	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	return f, nil
}
func WriteFile(f *os.File, text []byte) error {
	dataWriter := bufio.NewWriter(f)
	_, err := dataWriter.Write(text)
	if err != nil {
		return err
	}
	err = dataWriter.Flush()
	if err != nil {
		return err
	}
	return nil
}

func UpdateFile(t *testing.T, path string, data []byte) {
	f, err := CreateFile(path)
	require.NoError(t, err)
	err = WriteFile(f, data)
	require.NoError(t, err)
	err = f.Close()
	require.NoError(t, err)
}

func ReadFile(filePath string) ([]byte, error) {
	f, err := OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	err = f.Close()
	if err != nil {
		return nil, err
	}
	return content, err
}

func DoubleKeyIsExist(pos1 token.Pos, pos2 token.Pos, dict map[token.Pos]map[token.Pos]struct{}) bool {
	foundRaceA, isAExist := dict[pos1]
	if !isAExist {
		return false
	}
	_, isBExist := foundRaceA[pos2]
	if !isBExist {
		return false
	}
	return true
}
