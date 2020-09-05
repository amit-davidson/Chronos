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
	value   ssa.Value
	opKind  opKind
	lockset *lockset
}

type guardedAccessJSON struct {
	Value   string
	OpKind  opKind
	Lockset *lockset
}

func (ga *guardedAccess) MarshalJSON() ([]byte, error) {
	dumpJson := guardedAccessJSON{}
	dumpJson.Value = ga.value.Name()
	dumpJson.OpKind = ga.opKind
	dumpJson.Lockset = ga.lockset
	dump, err := json.Marshal(dumpJson)
	return dump, err
}
