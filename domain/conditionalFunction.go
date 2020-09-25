package domain

import "golang.org/x/tools/go/ssa"

type ConditionalFunction struct {
	Function      *ssa.CallCommon
	IsConditional bool
}
