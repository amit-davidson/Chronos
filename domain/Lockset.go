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

// The three following methods handle updating the lockset the same way. By recording each lock state (locked/unlocked)
// at the current point. It means that if a mutex was unlocked at some point but later was locked again, then it's latest
// status is locked, and the unlock status is removed. The difference between each algorithm is the context used.
//
// UpdateWithNewLockSet updates the lockset and expects newLocks, newUnlocks to contain the most up to date status of the
// mutex and update accordingly.
func (ls *Lockset) UpdateWithNewLockSet(newLocks, newUnlocks locksLastUse) {
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

// UpdateWithPrevLockset works the same as UpdateWithNewLockSet but expects prevLS to contain an earlier version of the
// status of the locks.
func (ls *Lockset) UpdateWithPrevLockset(prevLS *Lockset) {
	for lockName, lock := range prevLS.Locks {
		_, okLock := ls.Locks[lockName] // We check to see the lock doesn't exist to not override it with old reference of this lock
		_, okUnlock := ls.Unlocks[lockName]
		if !okLock && !okUnlock {
			ls.Locks[lockName] = lock
		}
	}

	for lockName, lock := range prevLS.Unlocks {
		_, okLock := ls.Locks[lockName] // We check to see the unlock doesn't exist to not override it with old reference of this unlock
		_, okUnlock := ls.Unlocks[lockName]
		if !okLock && !okUnlock {
			ls.Unlocks[lockName] = lock
		}
	}
}

// MergeSiblingLockset is called when merging different paths of the control flow graph. The mutex status should be
// merged and not appended. Because Locks is a must set, for a lock to appear in the result, an intersect
// between the branches' lockset is performed to make sure the lock appears in all branches. Unlock is a may set, so a
// union is applied since it's sufficient to have an unlock at least in one of the branches.
func (ls *Lockset) MergeSiblingLockset(locksetToMerge *Lockset) {
	locks := intersect(ls.Locks, locksetToMerge.Locks)
	unlocks := union(ls.Unlocks, locksetToMerge.Unlocks)

	for unlockName := range unlocks {
		// If there's a lock in one branch and an unlock in second, then unlock wins
		delete(locks, unlockName)
	}
	ls.Locks = locks
	ls.Unlocks = unlocks
}

func (ls *Lockset) Copy() *Lockset {
	newLs := &Lockset{}
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

func intersect(mapA, mapB locksLastUse) locksLastUse {
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

func union(mapA, mapB locksLastUse) locksLastUse {
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
