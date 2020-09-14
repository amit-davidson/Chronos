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

func NewEmptyLockSet() *Lockset {
	return &Lockset{
		ExistingLocks:   make(map[string]*ssa.CallCommon, 0),
		ExistingUnlocks: make(map[string]*ssa.CallCommon, 0),
	}
}

func NewLockSet(locks, unlocks map[string]*ssa.CallCommon) *Lockset {
	return &Lockset{
		ExistingLocks:   locks,
		ExistingUnlocks: unlocks,
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

func (ls *Lockset) Copy() *Lockset {
	newLs := NewEmptyLockSet()
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

func (ls *Lockset) ToJson() *LocksetJson {
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
	dump, err := json.Marshal(ls.ToJson())
	return dump, err
}
