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
	*PosData
	*FlowData
}

type PosData struct {
	PosID  int // guarded accesses of the same function share the same PosID. It's used to mark the same guarded access in different flows.
	Pos    token.Pos
	OpKind OpKind
	Value  ssa.Value
}

type FlowData struct {
	PosToRemove int
	ID          int // ID depends on the flow, which means it's unique.
	State       *Context
	Lockset     *Lockset
}

func (ga *FlowData) Copy() *FlowData {
	return &FlowData{
		ID:          ga.ID,
		PosToRemove: ga.PosToRemove,
		Lockset:     ga.Lockset.Copy(),
		State:       ga.State.CopyWithoutMap(),
	}
}

func (ga *GuardedAccess) Copy() *GuardedAccess {
	return &GuardedAccess{
		PosData:  ga.PosData,
		FlowData: ga.FlowData.Copy(),
	}
}

func (ga *GuardedAccess) ShallowCopy() *GuardedAccess {
	return &GuardedAccess{
		PosData:  ga.PosData,
		FlowData: ga.FlowData.Copy(),
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

func (ga *GuardedAccess) IsConflicting(gaToCompare *GuardedAccess) bool {
	return !ga.Intersects(gaToCompare) && ga.State.MayConcurrent(gaToCompare.State)
}

func AddGuardedAccess(pos token.Pos, value ssa.Value, kind OpKind, lockset *Lockset, context *Context) *GuardedAccess {
	context.Increment()
	return &GuardedAccess{
		PosData: &PosData{
			PosID:  PosIDCounter.GetNext(),
			Pos:    pos,
			OpKind: kind,
			Value:  value,
		},
		FlowData: &FlowData{
			ID:      GuardedAccessCounter.GetNext(),
			Lockset: lockset.Copy(),
			State:   context.CopyWithoutMap(),
		},
	}
}
