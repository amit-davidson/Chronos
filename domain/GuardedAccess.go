package domain

import (
	"encoding/json"
	"golang.org/x/tools/go/ssa"
)

type OpKind int

const (
	GuardAccessRead OpKind = iota
	GuardAccessWrite
)

type GuardedAccess struct {
	ID     int
	Value  ssa.Value
	State  *GoroutineState
	OpKind OpKind
}

type GuardedAccessJSON struct {
	Value  int
	OpKind OpKind
	State  *GoroutineStateJSON
}

func (ga *GuardedAccess) ToJSON() GuardedAccessJSON {
	dumpJson := GuardedAccessJSON{}
	dumpJson.Value = int(ga.Value.Pos())
	dumpJson.OpKind = ga.OpKind
	dumpJson.State = ga.State.ToJSON()
	return dumpJson
}
func (ga *GuardedAccess) MarshalJSON() ([]byte, error) {
	dump, err := json.Marshal(ga.ToJSON())
	return dump, err
}

func (ga *GuardedAccess) Intersects(gaToCompare *GuardedAccess) bool {
	if ga.ID == gaToCompare.ID || ga.State.GoroutineID == gaToCompare.State.GoroutineID {
		return true
	}
	if ga.OpKind == GuardAccessWrite || gaToCompare.OpKind == GuardAccessWrite {
		for _, lockA := range ga.State.Lockset.ExistingLocks {
			for _, lockB := range gaToCompare.State.Lockset.ExistingLocks {
				if lockA.Pos() == lockB.Pos() {
					return true
				}
			}
		}
	}
	return false
}
