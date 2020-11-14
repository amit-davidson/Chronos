package ssaPureUtils

import (
	"go/types"
	"golang.org/x/tools/go/ssa"
)

func GetField(value ssa.Value) (*ssa.FieldAddr, bool) {
	fieldAddr, ok := value.(*ssa.FieldAddr)
	return fieldAddr, ok
}

func GetUnderlyingObjectFromField(fieldAddr *ssa.FieldAddr) *types.Var {
	return fieldAddr.X.Type().Underlying().(*types.Pointer).Elem().Underlying().(*types.Struct).Field(fieldAddr.Field)
}

func FilterStructs(valueA, valueB ssa.Value) bool {
	fieldAddrA, okA := GetField(valueA)
	fieldAddrB, okB := GetField(valueB)

	isBothField := okA && okB
	if isBothField {
		fieldA := GetUnderlyingObjectFromField(fieldAddrA)
		fieldB := GetUnderlyingObjectFromField(fieldAddrB)
		if fieldA != fieldB { // If same struct but different fields
			return true
		}
	}

	return false
}
