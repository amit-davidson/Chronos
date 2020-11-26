package ssaPureUtils

import (
	"go/token"
	"golang.org/x/tools/go/ssa"
)


func GetMutexPos(value ssa.Value) token.Pos {
	val, ok := GetField(value)
	if !ok {
		return value.Pos()
	}
	obj := GetUnderlyingObjectFromField(val)
	return obj.Pos()
}