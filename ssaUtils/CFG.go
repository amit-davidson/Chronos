package ssaUtils

import (
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/utils/stacks"
	"golang.org/x/tools/go/ssa"
)

type CFG struct {
	visitedBlocksStack *stacks.BasicBlockStack
	ComputedBlocks     map[int]*domain.FunctionState
	calculatedState    *domain.FunctionState
}

func newCFG() *CFG {
	return &CFG{
		visitedBlocksStack: stacks.NewBasicBlockStack(),
		ComputedBlocks: make(map[int]*domain.FunctionState),
	}
}

func (cfg *CFG) CalculatePath(Context *domain.Context) {
	path := cfg.visitedBlocksStack.GetAllItems()
	block := path[0]
	state := cfg.ComputedBlocks[block.Index].Copy()
	for _, nextBlock := range path[1:] {
		nextState := cfg.ComputedBlocks[nextBlock.Index].Copy()
		for _, guardedAccess := range nextState.GuardedAccesses {
			guardedAccess.Lockset.UpdateLockSet(state.Lockset.ExistingLocks, state.Lockset.ExistingUnlocks)
		}
		state.MergeStates(nextState, true)
	}

	defersState := cfg.runDefers(Context, state.DeferredFunctions)
	state.MergeStates(defersState, true)

	if cfg.calculatedState == nil {
		cfg.calculatedState = state
	} else {
		cfg.calculatedState.MergeBranchState(state)
	}
}

func (cfg *CFG) calculateState(Context *domain.Context, block *ssa.BasicBlock) {
	cfg.visitedBlocksStack.Push(block)
	cfg.traverseGraph(Context, block)
	cfg.visitedBlocksStack.Pop()
}

func (cfg *CFG) traverseGraph(Context *domain.Context, block *ssa.BasicBlock) {
	if _, ok := cfg.ComputedBlocks[block.Index]; !ok {
		cfg.ComputedBlocks[block.Index] = GetBlockSummary(Context, block)
	}
	nextBlocks := block.Succs
	if len(nextBlocks) == 0 {
		cfg.CalculatePath(Context)
	} else {
		for _, nextBlock := range nextBlocks {
			cfg.visitedBlocksStack.Push(nextBlock)
			cfg.traverseGraph(Context, nextBlock)
			cfg.visitedBlocksStack.Pop()
		}
	}
}
