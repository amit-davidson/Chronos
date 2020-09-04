package main

import "golang.org/x/tools/go/ssa"

type opKind int

const (
	read opKind = iota
	write
)

type guardedAccess struct {
	value   ssa.Value
	opKind  opKind
	lockset lockset
}
