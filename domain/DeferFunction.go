package domain

import "golang.org/x/tools/go/ssa"

type DeferFunction struct {
	Function   *ssa.CallCommon
	BlockIndex int
}