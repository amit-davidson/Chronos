package main

import "golang.org/x/tools/go/ssa"

type lockset struct {
	existingLocks   map[string]*ssa.CallCommon
	existingUnlocks map[string]*ssa.CallCommon
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

func (ls *lockset) AddCallCommon(callCommon *ssa.CallCommon, isLocks bool) {
	receiver := callCommon.Args[0].(*ssa.Alloc).Comment
	locks := map[string]*ssa.CallCommon{receiver: callCommon}
	if isLocks {
		ls.updateLockSet(locks, nil)
	} else {
		ls.updateLockSet(nil, locks)
	}
}
