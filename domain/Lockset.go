package domain

import (
	"go/token"

	"golang.org/x/tools/go/ssa"
)

type locksLasUse map[token.Pos]*ssa.CallCommon

type Lockset struct {
	ExistingLocks   locksLasUse
	ExistingUnlocks locksLasUse
}

func NewLockset() *Lockset {
	return &Lockset{
		ExistingLocks:   make(locksLasUse, 0),
		ExistingUnlocks: make(locksLasUse, 0),
	}
}

func (ls *Lockset) UpdateLockSet(newLocks, newUnlocks locksLasUse) {
	for unlockName, _ := range newUnlocks {
		delete(ls.ExistingLocks, unlockName)
	}

	if newUnlocks != nil {
		for unlockName, unlock := range newUnlocks {
			ls.ExistingUnlocks[unlockName] = unlock
		}
	}
	for lockName, _ := range newLocks {
		if _, ok := ls.ExistingLocks[lockName]; ok {
			delete(ls.ExistingUnlocks, lockName)
		}
	}
}

func (ls *Lockset) MergeBranchesLockset(locksetToMerge *Lockset) {
	locks := Intersect(ls.ExistingLocks, locksetToMerge.ExistingLocks)
	unlocks := Union(ls.ExistingUnlocks, locksetToMerge.ExistingUnlocks)

	for unlockName, _ := range unlocks { // If there's a lock in one branch and an unlock in second, then unlock wins
		if _, ok := locks[unlockName]; ok {
			delete(locks, unlockName)
		}
	}
	ls.ExistingLocks = locks
	ls.ExistingUnlocks = unlocks
}

func (ls *Lockset) Copy() *Lockset {
	newLs := NewLockset()
	newLocks := make(locksLasUse)
	for key, value := range ls.ExistingLocks {
		newLocks[key] = value
	}
	newLs.ExistingLocks = newLocks
	newUnlocks := make(locksLasUse)
	for key, value := range ls.ExistingUnlocks {
		newUnlocks[key] = value
	}
	newLs.ExistingUnlocks = newUnlocks
	return newLs
}

func Intersect(mapA, mapB locksLasUse) locksLasUse {
	i := make(locksLasUse)
	for a := range mapA {
		for b := range mapB {
			if a == b {
				i[a] = mapA[a]
			}
		}
	}
	return i
}

func Union(mapA, mapB locksLasUse) locksLasUse {
	i := make(locksLasUse)
	for a := range mapA {
		i[a] = mapA[a]
	}
	for b := range mapB {
		i[b] = mapB[b]
	}
	return i
}
