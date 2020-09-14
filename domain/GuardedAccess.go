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

type guardedAccessJSON struct {
	Value       string
	OpKind      OpKind
	Lockset     *Lockset
}

func (ga *GuardedAccess) MarshalJSON() ([]byte, error) {
	dumpJson := guardedAccessJSON{}
	dumpJson.Value = ga.Value.Name()
	dumpJson.OpKind = ga.OpKind
	dumpJson.Lockset = ga.Lockset
	dump, err := json.Marshal(dumpJson)
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
