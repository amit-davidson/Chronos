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
	ID          string
	Value       ssa.Value
	OpKind      OpKind
	Lockset     *Lockset
	GoroutineId string
}

type GuardedAccessJSON struct {
	Value       int
	OpKind      OpKind
	Lockset     *LocksetJson
	GoroutineId string
}

func (ga *GuardedAccess) ToJson() GuardedAccessJSON {
	dumpJson := GuardedAccessJSON{}
	dumpJson.Value = int(ga.Value.Pos())
	dumpJson.OpKind = ga.OpKind
	dumpJson.GoroutineId = ga.GoroutineId
	dumpJson.Lockset = ga.Lockset.ToJson()
	return dumpJson
}
func (ga *GuardedAccess) MarshalJSON() ([]byte, error) {
	dump, err := json.Marshal(ga.ToJson())
	return dump, err
}

func (ga *GuardedAccess) Intersects(gaToCompare *GuardedAccess) bool {
	if ga.ID == gaToCompare.ID || ga.GoroutineId == gaToCompare.GoroutineId {
		return true
	}
	if ga.OpKind == GuardAccessWrite || gaToCompare.OpKind == GuardAccessWrite {
		for _, lockA := range ga.Lockset.ExistingLocks {
			for _, lockB := range gaToCompare.Lockset.ExistingLocks {
				if lockA.Pos() == lockB.Pos() {
					return true
				}
			}
		}
	}
	return false
}
