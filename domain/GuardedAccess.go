package domain

import (
	"github.com/amit-davidson/Chronos/ssaPureUtils"
	"go/token"
	"golang.org/x/tools/go/ssa"
)

type OpKind int

const (
	GuardAccessRead OpKind = iota
	GuardAccessWrite
)

func (op OpKind) String() string {
	switch op {
	case GuardAccessRead:
		return "Read"
	case GuardAccessWrite:
		return "Write"
	default:
		return "Unknown op type"
	}
}

type GuardedAccess struct {
	ID          int // ID depends on the flow, which means it's unique.
	PosID       int // guarded accesses of the same function share the same PosID. It's used to mark the same guarded access in different flows.
	Pos         token.Pos
	PosToRemove int
	Value       ssa.Value
	State       *Context
	Lockset     *Lockset
	OpKind      OpKind
}

func (ga *GuardedAccess) Copy() *GuardedAccess {
	return &GuardedAccess{
		ID:          ga.ID,
		PosID:       ga.PosID,
		Pos:         ga.Pos,
		PosToRemove: ga.PosToRemove,
		Value:       ga.Value,
		Lockset:     ga.Lockset.Copy(),
		OpKind:      ga.OpKind,
		State:       ga.State.Copy(),
	}
}

func (ga *GuardedAccess) ShallowCopy() *GuardedAccess {
	return &GuardedAccess{
		ID:          ga.ID,
		PosID:       ga.PosID,
		Pos:         ga.Pos,
		PosToRemove: ga.PosToRemove,
		Value:       ga.Value,
		Lockset:     ga.Lockset,
		OpKind:      ga.OpKind,
		State:       ga.State,
	}
}

func (ga *GuardedAccess) Intersects(gaToCompare *GuardedAccess) bool {
	if ga.ID == gaToCompare.ID || ga.State.GoroutineID == gaToCompare.State.GoroutineID {
		return true
	}
	if ga.OpKind == GuardAccessRead && gaToCompare.OpKind == GuardAccessRead {
		return true
	}

	if ssaPureUtils.FilterStructs(ga.Value, gaToCompare.Value) {
		return true
	}

	for lockA := range ga.Lockset.Locks {
		for lockB := range gaToCompare.Lockset.Locks {
			if lockA == lockB {
				return true
			}
		}
	}
	return false
}

func AddGuardedAccess(pos token.Pos, value ssa.Value, kind OpKind, lockset *Lockset, context *Context) *GuardedAccess {
	context.Increment()
	return &GuardedAccess{ID: GuardedAccessCounter.GetNext(), PosID: PosIDCounter.GetNext(), Pos: pos,
		Value: value, Lockset: lockset.Copy(), OpKind: kind, State: context.Copy()}
}
