package testutils

import (
	"encoding/json"
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/ssaUtils"
	"go/token"
)

type GuardedAccessJSON struct {
	ID          int
	Pos         string
	Value       string
	OpKind      domain.OpKind
	LocksetJson *LocksetJson
	State       *ContextJSON
}

func GuardedAccessToJSON(ga *domain.GuardedAccess) GuardedAccessJSON {
	prog := ssaUtils.GlobalProgram
	dumpJson := GuardedAccessJSON{}
	dumpJson.ID = ga.ID
	dumpJson.Pos = prog.Fset.Position(ga.Pos).String()
	dumpJson.Value = prog.Fset.Position(ga.Value.Pos()).String()
	dumpJson.OpKind = ga.OpKind
	dumpJson.LocksetJson = LocksetToJSON(ga.Lockset)
	dumpJson.State = ContextToJSON(ga.State)
	return dumpJson
}
func MarshalJSON(ga *domain.GuardedAccess) ([]byte, error) {
	dump, err := json.Marshal(GuardedAccessToJSON(ga))
	return dump, err
}

// Lockset name ans pos
type LocksetJson struct {
	ExistingLocks   map[string]token.Position
	ExistingUnlocks map[string]token.Position
}

func LocksetToJSON(ls *domain.Lockset) *LocksetJson {
	prog := ssaUtils.GlobalProgram

	dumpJson := &LocksetJson{ExistingLocks: make(map[string]token.Position, 0), ExistingUnlocks: make(map[string]token.Position, 0)}
	for lockName, lock := range ls.ExistingLocks {
		lockPos := prog.Fset.Position(lock.Pos())
		dumpJson.ExistingLocks[lockName] = lockPos
	}
	for lockName, lock := range ls.ExistingUnlocks {
		unlockPos := prog.Fset.Position(lock.Pos())
		dumpJson.ExistingUnlocks[lockName] = unlockPos
	}
	return dumpJson
}
func MarshalLocksetJSON(ls *domain.Lockset) ([]byte, error) {
	dump, err := json.Marshal(LocksetToJSON(ls))
	return dump, err
}

type ContextJSON struct {
	GoroutineID int
	Clock       domain.VectorClock
}

func ContextToJSON(gs *domain.Context) *ContextJSON {
	dumpJson := ContextJSON{}
	dumpJson.GoroutineID = gs.GoroutineID
	dumpJson.Clock = gs.Clock
	return &dumpJson
}
