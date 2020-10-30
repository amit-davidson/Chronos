package domain

import (
	"go/token"

	"golang.org/x/tools/go/ssa"
)

type locksLastUse map[token.Pos]*ssa.CallCommon

type Lockset struct {
	ExistingLocks   locksLastUse
	ExistingUnlocks locksLastUse
}

func NewLockset() *Lockset {
	return &Lockset{
		ExistingLocks:   make(locksLastUse),
		ExistingUnlocks: make(locksLastUse),
	}
}

func (ls *Lockset) UpdateLockSet(newLocks, newUnlocks locksLastUse) {
	// The algorithm works by remembering each lock's state (locked/unlocked/or nothing, of course).
	// It means that if a mutex was unlocked at some point but later was locked again,
	// then its latest status is locked, and the unlock status is removed.
	// Source: https://github.com/amit-davidson/Chronos/pull/10/files#r507203577
	for lockName, lock := range newLocks {
		ls.ExistingLocks[lockName] = lock
	}
	for unlockName := range newUnlocks {
		delete(ls.ExistingLocks, unlockName)
	}
	for unlockName, unlock := range newUnlocks {
		ls.ExistingUnlocks[unlockName] = unlock
	}
	for lockName := range newLocks {
		if _, ok := ls.ExistingLocks[lockName]; ok {
			delete(ls.ExistingUnlocks, lockName)
		}
	}
}

func (ls *Lockset) MergeBranchesLockset(locksetToMerge *Lockset) {
	locks := Intersect(ls.ExistingLocks, locksetToMerge.ExistingLocks)
	unlocks := Union(ls.ExistingUnlocks, locksetToMerge.ExistingUnlocks)

	for unlockName := range unlocks {
		// If there's a lock in one branch and an unlock in second, then unlock wins
		delete(locks, unlockName)
	}
	ls.ExistingLocks = locks
	ls.ExistingUnlocks = unlocks
}

func (ls *Lockset) Copy() *Lockset {
	newLs := NewLockset()
	newLocks := make(locksLastUse, len(ls.ExistingLocks))
	for key, value := range ls.ExistingLocks {
		newLocks[key] = value
	}
	newLs.ExistingLocks = newLocks
	newUnlocks := make(locksLastUse, len(ls.ExistingUnlocks))
	for key, value := range ls.ExistingUnlocks {
		newUnlocks[key] = value
	}
	newLs.ExistingUnlocks = newUnlocks
	return newLs
}

func Intersect(mapA, mapB locksLastUse) locksLastUse {
	i := make(locksLastUse, min(len(mapA), len(mapB)))
	for a := range mapA {
		for b := range mapB {
			if a == b {
				i[a] = mapA[a]
			}
		}
	}
	return i
}

func Union(mapA, mapB locksLastUse) locksLastUse {
	i := make(locksLastUse, max(len(mapA), len(mapB)))
	for a := range mapA {
		i[a] = mapA[a]
	}
	for b := range mapB {
		i[b] = mapB[b]
	}
	return i
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
