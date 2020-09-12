package main

import (
	"encoding/json"
	"golang.org/x/tools/go/ssa"
)

type opKind int

const (
	read opKind = iota
	write
)

type guardedAccess struct {
	id          string
	value       ssa.Value
	opKind      opKind
	lockset     *lockset
	GoroutineId string
}

type guardedAccessJSON struct {
	Value       string
	OpKind      opKind
	Lockset     *lockset
	GoroutineId string
}

func (ga *guardedAccess) MarshalJSON() ([]byte, error) {
	dumpJson := guardedAccessJSON{}
	dumpJson.Value = ga.value.Name()
	dumpJson.OpKind = ga.opKind
	dumpJson.Lockset = ga.lockset
	dumpJson.GoroutineId = ga.GoroutineId
	dump, err := json.Marshal(dumpJson)
	return dump, err
}

func (ga *guardedAccess) Intersects(gaToCompare *guardedAccess) bool {
	if ga.id == gaToCompare.id || ga.GoroutineId == gaToCompare.GoroutineId {
		return true
	}
	if ga.opKind == write || gaToCompare.opKind == write {
		for _, lockA := range ga.lockset.existingLocks {
			for _, lockB := range gaToCompare.lockset.existingLocks {
				if lockA.Pos() == lockB.Pos() {
					return true
				}
			}
		}
	}
	return false
}
