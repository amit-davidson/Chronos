package utils

import (
	"bufio"
	"github.com/stretchr/testify/require"
	"go/types"
	"golang.org/x/tools/go/ssa"
	"io/ioutil"
	"os"
	"testing"
)

func IsCallTo(call *ssa.Function, names ...string) bool {
	fn, ok := call.Object().(*types.Func)
	if !ok {
		return false
	}
	q := fn.FullName()
	for _, name := range names {
		if q == name {
			return true
		}
	}
	return false
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

func ReadLineByNumber(filePath string, lineNumber int) (string, error) {
	f, err := OpenFile(filePath)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(f)
	line := ""
	for i := 1; scanner.Scan(); i++ {
		if i == lineNumber {
			line = scanner.Text()
			break
		}
	}

	err = f.Close()
	if err != nil {
		return "", err
	}
	return line, err
}
