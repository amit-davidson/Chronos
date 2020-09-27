package domain

import "golang.org/x/tools/go/ssa"

type FunctionWithBlock struct {
	Function   *ssa.CallCommon
	BlockIndex int
}