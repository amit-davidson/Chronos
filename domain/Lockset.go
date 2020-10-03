package domain

import (
	"encoding/json"
	"golang.org/x/tools/go/ssa"
)

type Lockset struct {
	ExistingLocks   map[string]*ssa.CallCommon
	ExistingUnlocks map[string]*ssa.CallCommon
}

// Lockset name ans pos
type LocksetJson struct {
	ExistingLocks   map[string]int
	ExistingUnlocks map[string]int
}

func NewLockset() *Lockset {
	return &Lockset{
		ExistingLocks:   make(map[string]*ssa.CallCommon, 0),
		ExistingUnlocks: make(map[string]*ssa.CallCommon, 0),
	}
}

func (ls *Lockset) UpdateLockSet(newLocks, newUnlocks map[string]*ssa.CallCommon) {
	if newLocks != nil {
		for lockName, lock := range newLocks {
			ls.ExistingLocks[lockName] = lock
		}
	}
	for unlockName, _ := range newUnlocks {
		if _, ok := ls.ExistingLocks[unlockName]; ok {
			delete(ls.ExistingLocks, unlockName)
		}
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
	newLocks := make(map[string]*ssa.CallCommon)
	for key, value := range ls.ExistingLocks {
		newLocks[key] = value
	}
	newLs.ExistingLocks = newLocks
	newUnlocks := make(map[string]*ssa.CallCommon)
	for key, value := range ls.ExistingUnlocks {
		newUnlocks[key] = value
	}
	newLs.ExistingUnlocks = newUnlocks
	return newLs
}

func (ls *Lockset) ToJSON() *LocksetJson {
	dumpJson := &LocksetJson{ExistingLocks: make(map[string]int, 0), ExistingUnlocks: make(map[string]int, 0)}
	for lockName, lock := range ls.ExistingLocks {
		dumpJson.ExistingLocks[lockName] = int(lock.Pos())
	}
	for lockName, lock := range ls.ExistingUnlocks {
		dumpJson.ExistingUnlocks[lockName] = int(lock.Pos())
	}
	return dumpJson
}
func (ls *Lockset) MarshalJSON() ([]byte, error) {
	dump, err := json.Marshal(ls.ToJSON())
	return dump, err
}

func Intersect(mapA, mapB map[string]*ssa.CallCommon) map[string]*ssa.CallCommon {
	i := make(map[string]*ssa.CallCommon)
	for a := range mapA {
		for b := range mapB {
			if a == b {
				i[a] = mapA[a]
			}
		}
	}
	return i
}

func Union(mapA, mapB map[string]*ssa.CallCommon) map[string]*ssa.CallCommon {
	i := make(map[string]*ssa.CallCommon)
	for a := range mapA {
		i[a] = mapA[a]
	}
	for b := range mapB {
		i[b] = mapB[b]
	}
	return i
}