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
)

// Test of context-sensitive treatment of certain function calls,
// e.g. static calls to simple accessor methods.

var a, b int

type T struct{ x *int }

func (t *T) SetX(x *int) { t.x = x }
func (t *T) GetX() *int  { return t.x }

func context3() {
	mutex := sync.Mutex{}
	mutex2 := sync.Mutex{}
	mutex3 := sync.Mutex{}
	a := false
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
	defer mutex.Lock()
	if a {
		mutex.Unlock()
	}
	defer mutex.Unlock()
	map1 := map[string]map[string]interface{}{}
	map1["map2"] = map[string]interface{}{}
	map1["map2"]["map2"] = map[string]interface{}{}
}

func main() {
	context3()
}
`

func main() {
	var conf loader.Config
	// Parse the input file, a string.
	// (Command-line tools should use conf.FromArgs.)
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

	// Configure the pointer analysis to build a call-graph.
	//config := &pointer.Config{
	//	Mains:          []*ssa.Package{mainPkg},
	//	BuildCallGraph: true,
	//}
	//
	//calls := make(map[*ssa.CallCommon]bool)
	// Query points-to set of (C).f's parameter m, a map.
	funcInterface5 := mainPkg.Func("context3")
	_ = GetFunctionLocks(funcInterface5)
	//C := mainPkg.Type("P").Type()
	//funcInterface5 := prog.LookupMethod(C, mainPkg.Pkg, "f")
	//for _, b := range funcInterface5.Blocks {
	//	for _, instr := range b.Instrs {
	//		if instr, ok := instr.(ssa.CallInstruction); ok {
	//			call := instr.Common()
	//			if b, ok := call.Value.(*ssa.Builtin); ok && b.Name() == "print" && len(call.Args) == 1 {
	//				calls[instr.Common()] = true
	//			}
	//		}
	//	}
	//}
	//for probe := range calls {
	//	v := probe.Args[0]
	//	if pointer.CanPoint(v.Type()) {
	//		config.AddQuery(v)
	//	}
	//}
	//
	////Cfm := funcInterface5.Params[1]
	////config.AddQuery(Cfm)
	//
	//// Run the pointer analysis.
	//result, err := pointer.Analyze(config)
	//if err != nil {
	//	panic(err) // internal error in pointer analysis
	//}
	//
	//// Find edges originating from the main package.
	//// By converting to strings, we de-duplicate nodes
	//// representing the same function due to context sensitivity.
	//var edges []string
	//callgraph.GraphVisitEdges(result.CallGraph, func(edge *callgraph.Edge) error {
	//	caller := edge.Caller.Func
	//	if caller.Pkg == mainPkg {
	//		edges = append(edges, fmt.Sprint(caller, " --> ", edge.Callee.Func))
	//	}
	//	return nil
	//})
	//
	//// Print the edges in sorted order.
	//sort.Strings(edges)
	//for _, edge := range edges {
	//	fmt.Println(edge)
	//}
	//fmt.Println()
	//
	//// Print the labels of (C).f(m)'s points-to set.
	//fmt.Println("vars:")
	//var labels []string
	////pts := result.Queries[Cfm].PointsTo()
	//for query, queryRes := range result.Queries {
	//	_, ok := query.(*ssa.Extract)
	//	if ok {
	//		queryPos := prog.Fset.Position(query.(*ssa.Extract).Tuple.Pos())
	//		fmt.Println(fmt.Sprintf(" query: %s", queryPos))
	//	}
	//	pts := queryRes.PointsTo()
	//	for _, l := range pts.Labels() {
	//		label := fmt.Sprintf("  %s: %s", prog.Fset.Position(l.Pos()), l)
	//		labels = append(labels, label)
	//	}
	//	sort.Strings(labels)
	//	for _, label := range labels {
	//		fmt.Println(label)
	//	}
	//}
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

func GetBlockLocks(block *ssa.BasicBlock) (*lockset, []*lockset) {
	ls := newEmptyLockSet()
	deferredCalls := make([]*lockset, 0)

	instrs := FilterDebug(block.Instrs)
	for _, ins := range instrs[:len(instrs)-1] {
		switch call := ins.(type) {
		case *ssa.Call:
			callCommon := call.Common()
			receiver := callCommon.Args[0].(*ssa.Alloc).Comment
			locks := map[string]*ssa.CallCommon{receiver: callCommon}
			if IsCallToAny(callCommon, "(*sync.Mutex).Lock") {
				ls.updateLockSet(locks, nil)
			}
			if IsCallToAny(callCommon, "(*sync.Mutex).Unlock") {
				ls.updateLockSet(nil, locks)
			}
			continue

		case *ssa.Defer:
			callCommon := call.Common()
			receiver := callCommon.Args[0].(*ssa.Alloc).Comment
			locks := map[string]*ssa.CallCommon{receiver: callCommon}
			if IsCallToAny(callCommon, "(*sync.Mutex).Lock") {
				deferredCalls = append(deferredCalls, newLockSet(locks, nil))
			}
			if IsCallToAny(callCommon, "(*sync.Mutex).Unlock") {
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
		lsRet, deferredCallsRet := GetBlockLocks(block)
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
