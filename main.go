package main

import (
	"fmt"
	"go/types"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

const myprog = `
package main

import (
	"sync"
	"fmt"
    "math/rand"
)

// Test of context-sensitive treatment of certain function calls,
// e.g. static calls to simple accessor methods.

var a, b int
var t = rand.Int()
var cond = false
	//map1 := map[string]map[string]interface{}{}
type T struct{ X int }

func (t *T) SetX(x int) { t.X = x }
func (t *T) GetX() int  { return t.X }

var ls1 = make([]int, 5)
var ls2 = make([]int, 5, 5)

 //go:noinline
func context2 (a int) {
	fmt.Print("4")
}
func context3() {
	mutex := sync.Mutex{}
	mutex2 := sync.Mutex{}
	mutex3 := sync.Mutex{}
	s := 1
	_ = s
	a := false
	//tStruct := T{X:1}
	if a {
		defer mutex.Unlock()
	} else if a {
		mutex.Unlock()
	} else {
		mutex.Unlock()
		mutex2.Unlock()
	}
	mutex.Lock()
	mutex3.Lock()
	//defer mutex.Lock()
	//if a {
	//	mutex.Unlock()
	//}
	//defer mutex.Unlock()
	if cond {
		map1 := map[string]map[string]int{}
		map1["map2"] = map[string]int{}
		//map1["map2"]["map2"] = map[string]interface{}{}
		map1["map2"]["map2"] = 5
		//ls1[0] = 3
		//ls2[1] = ls1[1]
		//tStruct.X = 2 + rand.Int()
		//context2(tStruct.X)
	}
}

func main() {
	context3()
}
`

func main() {
	var conf loader.Config
	file, err := conf.ParseFile("myprog.go", myprog)
	if err != nil {
		fmt.Print(err) // parse error
		return
	}

	// Create single-file main package and import its dependencies.
	conf.CreateFromFiles("main", file)

	iprog, err := conf.Load()
	if err != nil {
		fmt.Print(err) // type error in some package
		return
	}

	// Create SSA-form program representation.
	prog := ssautil.CreateProgram(iprog, 0)
	mainPkg := prog.Package(iprog.Created[0].Pkg)

	// Build SSA code for bodies of all functions in the whole program.
	prog.Build()

	funcInterface5 := mainPkg.Func("context3")
	_ = GetFunctionLocks(funcInterface5)
}

func IsCallToAny(call *ssa.CallCommon, names ...string) bool {
	q := CallName(call, false)
	for _, name := range names {
		if q == name {
			return true
		}
	}
	return false
}

func CallName(call *ssa.CallCommon, short bool) string {
	if call.IsInvoke() {
		return ""
	}
	switch v := call.Value.(type) {
	case *ssa.Function:
		fn, ok := v.Object().(*types.Func)
		if !ok {
			return ""
		}
		if short {
			return fn.Name()
		} else {
			return fn.FullName()
		}
	case *ssa.Builtin:
		return v.Name()
	}
	return ""
}

func FilterDebug(instr []ssa.Instruction) []ssa.Instruction {
	var out []ssa.Instruction
	for _, ins := range instr {
		if _, ok := ins.(*ssa.DebugRef); !ok {
			out = append(out, ins)
		}
	}
	return out
}

func addGuardedAccess(guardedAccesses []*guardedAccess, value ssa.Value, kind opKind, currentLockset lockset) {
	guardedAccessToAdd := &guardedAccess{value: value, opKind: kind, lockset: currentLockset}
	guardedAccesses = append(guardedAccesses, guardedAccessToAdd)
}

func GetBlockLocks(block *ssa.BasicBlock, ls *lockset) (*lockset, []*lockset) {
	currentLockset := ls
	deferredCalls := make([]*lockset, 0)
	guardedAccesses := make([]*guardedAccess, 0)
	instrs := FilterDebug(block.Instrs)
	for _, ins := range instrs[:len(instrs)-1] {
		switch call := ins.(type) {
		case *ssa.BinOp:
			addGuardedAccess(guardedAccesses, call.X, read, *currentLockset)
			addGuardedAccess(guardedAccesses, call.Y, read, *currentLockset)
		case *ssa.UnOp:
			addGuardedAccess(guardedAccesses, call.X, read, *currentLockset)
		case *ssa.Store:
			addGuardedAccess(guardedAccesses, call.Addr, write, *currentLockset)
			addGuardedAccess(guardedAccesses, call.Val, read, *currentLockset)
		case *ssa.Field:
			addGuardedAccess(guardedAccesses, call.X, read, *currentLockset)
		case *ssa.FieldAddr:
			addGuardedAccess(guardedAccesses, call.X, read, *currentLockset)
		case *ssa.Index:
			addGuardedAccess(guardedAccesses, call.X, read, *currentLockset)
		case *ssa.IndexAddr:
			addGuardedAccess(guardedAccesses, call.X, read, *currentLockset)
		case *ssa.Lookup:
			addGuardedAccess(guardedAccesses, call.X, read, *currentLockset)
		case *ssa.MapUpdate:
			addGuardedAccess(guardedAccesses, call.Map, write, *currentLockset)
			addGuardedAccess(guardedAccesses, call.Value, read, *currentLockset)
		case *ssa.Panic:
			addGuardedAccess(guardedAccesses, call.X, read, *currentLockset)
		case *ssa.Call:
			callCommon := call.Common()
			if IsCallToAny(callCommon, "(*sync.Mutex).Lock") {
				receiver := callCommon.Args[0].(*ssa.Alloc).Comment
				locks := map[string]*ssa.CallCommon{receiver: callCommon}
				ls.updateLockSet(locks, nil)
			}
			if IsCallToAny(callCommon, "(*sync.Mutex).Unlock") {
				receiver := callCommon.Args[0].(*ssa.Alloc).Comment
				locks := map[string]*ssa.CallCommon{receiver: callCommon}
				ls.updateLockSet(nil, locks)
			}
			continue

		case *ssa.Defer:
			callCommon := call.Common()
			if IsCallToAny(callCommon, "(*sync.Mutex).Lock") {
				receiver := callCommon.Args[0].(*ssa.Alloc).Comment
				locks := map[string]*ssa.CallCommon{receiver: callCommon}
				deferredCalls = append(deferredCalls, newLockSet(locks, nil))
			}
			if IsCallToAny(callCommon, "(*sync.Mutex).Unlock") {
				receiver := callCommon.Args[0].(*ssa.Alloc).Comment
				locks := map[string]*ssa.CallCommon{receiver: callCommon}
				deferredCalls = append(deferredCalls, newLockSet(nil, locks))
			}
			continue
		}
	}
	return ls, deferredCalls
}

func GetFunctionLocks(fn *ssa.Function) *lockset {
	var conditionalBlocks = map[string]struct{}{
		"if.then": {},
		"if.else": {},
	}
	ls := newEmptyLockSet()
	deferredCalls := make([]*lockset, 0)

	for _, block := range fn.Blocks {
		lsRet, deferredCallsRet := GetBlockLocks(block, ls)
		if _, ok := conditionalBlocks[block.Comment]; ok {
			ls.updateLockSet(nil, lsRet.existingUnlocks) // Ignore locks in a condition branch since it's a must set.
			for _, deferredCallRet := range deferredCallsRet {
				deferredCalls = append(deferredCalls, newLockSet(nil, deferredCallRet.existingUnlocks))
			}
		} else {
			ls.updateLockSet(lsRet.existingLocks, lsRet.existingUnlocks)
			for _, deferredCallRet := range deferredCallsRet {
				deferredCalls = append(deferredCalls, newLockSet(deferredCallRet.existingLocks, deferredCallRet.existingUnlocks))
			}
		}
	}

	for i := len(deferredCalls) - 1; i >= 0; i-- {
		ls.updateLockSet(deferredCalls[i].existingLocks, deferredCalls[i].existingUnlocks)
	}
	return ls
}
