package ssaUtils

import (
	"github.com/amit-davidson/Chronos/domain"
	"go/token"
	"golang.org/x/tools/go/ssa"
)

func GetStackTrace(prog *ssa.Program, ga *domain.GuardedAccess) string {
	stack := ""
	for _, pos := range ga.State.StackTrace.Iter() {
		calculatedPos := prog.Fset.Position(token.Pos(pos))
		stack += calculatedPos.String()
		stack += " ->\n"
	}
	return stack
}
