package main

import (
	"encoding/json"
	"golang.org/x/tools/go/ssa"
)

type lockset struct {
	existingLocks   map[string]*ssa.CallCommon
	existingUnlocks map[string]*ssa.CallCommon
}

// lockset name ans pos
type locksetJson struct {
	ExistingLocks   map[string]int
	ExistingUnlocks map[string]int
}

func newEmptyLockSet() *lockset {
	return &lockset{
		existingLocks:   make(map[string]*ssa.CallCommon, 0),
		existingUnlocks: make(map[string]*ssa.CallCommon, 0),
	}
}

func newLockSet(locks, unlocks map[string]*ssa.CallCommon) *lockset {
	return &lockset{
		existingLocks:   locks,
		existingUnlocks: unlocks,
	}
}

func (ls *lockset) updateLockSet(newLocks, newUnlocks map[string]*ssa.CallCommon) {
	if newLocks != nil {
		for lockName, lock := range newLocks {
			ls.existingLocks[lockName] = lock
		}
	}
	for unlockName, _ := range newUnlocks {
		if _, ok := ls.existingLocks[unlockName]; ok {
			delete(ls.existingLocks, unlockName)
		}
	}

	if newUnlocks != nil {
		for unlockName, unlock := range newUnlocks {
			ls.existingUnlocks[unlockName] = unlock
		}
	}
	for lockName, _ := range newLocks {
		if _, ok := ls.existingLocks[lockName]; ok {
			delete(ls.existingUnlocks, lockName)
		}
	}
}

func (ls *lockset) Copy() *lockset {
	newLs := newEmptyLockSet()
	newLocks := make(map[string]*ssa.CallCommon)
	for key, value := range ls.existingLocks {
		newLocks[key] = value
	}
	newLs.existingLocks = newLocks
	newUnlocks := make(map[string]*ssa.CallCommon)
	for key, value := range ls.existingUnlocks {
		newUnlocks[key] = value
	}
	newLs.existingUnlocks = newUnlocks
	return newLs
}

func (ls *lockset) MarshalJSON() ([]byte, error) {
	dumpJson := locksetJson{ExistingLocks: make(map[string]int, 0), ExistingUnlocks: make(map[string]int, 0)}
	for lockName, lock := range ls.existingLocks {
		dumpJson.ExistingLocks[lockName] = int(lock.Pos())
	}
	for lockName, lock := range ls.existingUnlocks {
		dumpJson.ExistingUnlocks[lockName] = int(lock.Pos())
	}
	dump, err := json.Marshal(dumpJson)
	return dump, err
}
