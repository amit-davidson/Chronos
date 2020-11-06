package domain

import (
	"go/token"

	"golang.org/x/tools/go/ssa"
)

type locksLastUse map[token.Pos]*ssa.CallCommon

type Lockset struct {
	Locks   locksLastUse
	Unlocks locksLastUse
}

func NewLockset() *Lockset {
	return &Lockset{
		Locks:   make(locksLastUse),
		Unlocks: make(locksLastUse),
	}
}

func (ls *Lockset) UpdateLockSet(newLocks, newUnlocks locksLastUse) {
	// The algorithm works by remembering each lock's state (locked/unlocked/or nothing, of course).
	// It means that if a mutex was unlocked at some point but later was locked again,
	// then its latest status is locked, and the unlock status is removed.
	// Source: https://github.com/amit-davidson/Chronos/pull/10/files#r507203577
	for lockName, lock := range newLocks {
		ls.Locks[lockName] = lock
	}
	for unlockName := range newUnlocks {
		delete(ls.Locks, unlockName)
	}
	for unlockName, unlock := range newUnlocks {
		ls.Unlocks[unlockName] = unlock
	}
	for lockName := range newLocks {
		if _, ok := ls.Locks[lockName]; ok {
			delete(ls.Unlocks, lockName)
		}
	}
}

func (ls *Lockset) MergeSiblingLockset(locksetToMerge *Lockset) {
	locks := Intersect(ls.Locks, locksetToMerge.Locks)
	unlocks := Union(ls.Unlocks, locksetToMerge.Unlocks)

	for unlockName := range unlocks {
		// If there's a lock in one branch and an unlock in second, then unlock wins
		delete(locks, unlockName)
	}
	ls.Locks = locks
	ls.Unlocks = unlocks
}

func (ls *Lockset) Copy() *Lockset {
	newLs := NewLockset()
	newLocks := make(locksLastUse, len(ls.Locks))
	for key, value := range ls.Locks {
		newLocks[key] = value
	}
	newLs.Locks = newLocks
	newUnlocks := make(locksLastUse, len(ls.Unlocks))
	for key, value := range ls.Unlocks {
		newUnlocks[key] = value
	}
	newLs.Unlocks = newUnlocks
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
