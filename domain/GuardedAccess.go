package domain

import (
	"Chronos/utils"
	"encoding/json"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ssa"
)

type OpKind int

const (
	GuardAccessRead OpKind = iota
	GuardAccessWrite
)

func (op OpKind) String() string {
	switch op {
	case GuardAccessRead:
		return "Read"
	case GuardAccessWrite:
		return "Write"
	default:
		return "Unknown op type"
	}
}

type GuardedAccess struct {
	ID         int
	Pos        token.Pos
	Stacktrace []int
	Value      ssa.Value
	State      *Context
	Lockset    *Lockset
	OpKind     OpKind
}

type GuardedAccessJSON struct {
	ID          int
	Pos         token.Pos
	Value       int
	OpKind      OpKind
	LocksetJson *LocksetJson
	State       *ContextJSON
}

func (ga *GuardedAccess) ToJSON() GuardedAccessJSON {
	dumpJson := GuardedAccessJSON{}
	dumpJson.ID = ga.ID
	dumpJson.Pos = ga.Pos
	dumpJson.Value = int(ga.Value.Pos())
	dumpJson.OpKind = ga.OpKind
	dumpJson.LocksetJson = ga.Lockset.ToJSON()
	dumpJson.State = ga.State.ToJSON()
	return dumpJson
}
func (ga *GuardedAccess) MarshalJSON() ([]byte, error) {
	dump, err := json.Marshal(ga.ToJSON())
	return dump, err
}

func FilterStructs(valueA, valueB ssa.Value) bool {
	fieldAddrA, okA := valueA.(*ssa.FieldAddr)
	fieldAddrB, okB := valueB.(*ssa.FieldAddr)
	//isOnlyOneField := okA && !okB || okB && !okA
	//if isOnlyOneField { // If a field points to a struct
	//	return true
	//}

	isBothField := okA && okB
	if isBothField {
		fieldA := fieldAddrA.X.Type().Underlying().(*types.Pointer).Elem().Underlying().(*types.Struct).Field(fieldAddrA.Field)
		fieldB := fieldAddrB.X.Type().Underlying().(*types.Pointer).Elem().Underlying().(*types.Struct).Field(fieldAddrB.Field)
		if fieldA != fieldB { // If same struct but different fields
			return true
		}
	}

	return false
}

func (ga *GuardedAccess) Intersects(gaToCompare *GuardedAccess) bool {
	if ga.ID == gaToCompare.ID || ga.State.GoroutineID == gaToCompare.State.GoroutineID {
		return true
	}
	if ga.OpKind == GuardAccessRead && gaToCompare.OpKind == GuardAccessRead {
		return true
	}

	if FilterStructs(ga.Value, gaToCompare.Value) {
		return true
	}

	for lockA := range ga.Lockset.ExistingLocks {
		for lockB := range gaToCompare.Lockset.ExistingLocks {
			if lockA == lockB {
				return true
			}
		}
	}
	return false
}

var GuardedAccessCounter = utils.NewCounter()

func AddGuardedAccess(pos token.Pos, value ssa.Value, kind OpKind, lockset *Lockset, Context *Context) *GuardedAccess {
	Context.Increment()
	stackTrace := Context.StackTrace.GetAllItems()
	return &GuardedAccess{ID: GuardedAccessCounter.GetNext(), Pos: pos, Value: value, Lockset: lockset.Copy(), OpKind: kind, Stacktrace: stackTrace, State: Context.Copy()}
}
