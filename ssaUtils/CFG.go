package ssaUtils

import (
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/utils/stacks"
	"golang.org/x/tools/go/ssa"
)

type CFG struct {
	visitedBlocksStack *stacks.BasicBlockStack

	ComputedBlocks      map[int]*domain.FunctionState
	ComputedDeferBlocks map[int]*domain.FunctionState
	calculatedState     *domain.FunctionState
}

func newCFG() *CFG {
	return &CFG{
		visitedBlocksStack:  stacks.NewBasicBlockStack(),
		ComputedBlocks:      make(map[int]*domain.FunctionState),
		ComputedDeferBlocks: make(map[int]*domain.FunctionState),
	}
}

func (cfg *CFG) CalculatePath() {
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

	deferBlock := path[0]
	deferState := cfg.ComputedDeferBlocks[deferBlock.Index]
	if deferState != nil {
		deferStateCopy := deferState.Copy()
		for _, nextBlock := range path[1:] {
			nextState := cfg.ComputedDeferBlocks[nextBlock.Index]
			if nextState == nil {
				continue
			}
			nextStateCopy := nextState.Copy()
			for _, guardedAccess := range nextStateCopy.GuardedAccesses {
				guardedAccess.Lockset.UpdateLockSet(state.Lockset.ExistingLocks, state.Lockset.ExistingUnlocks)
			}
			deferStateCopy.MergeStates(nextStateCopy, true)
		}

		state.MergeStates(deferStateCopy, true)
	}
	if cfg.calculatedState == nil {
		cfg.calculatedState = state
	} else {
		cfg.calculatedState.MergeBranchState(state)
	}
}

func (cfg *CFG) calculateState(Context *domain.Context, block *ssa.BasicBlock) {
	firstBlock := &ssa.BasicBlock{Index: -1, Succs: []*ssa.BasicBlock{block}}
	cfg.traverseGraph(Context, firstBlock)
}

func (cfg *CFG) traverseGraph(Context *domain.Context, block *ssa.BasicBlock) {
	nextBlocks := block.Succs
	for _, nextBlock := range nextBlocks {
		if _, ok := cfg.ComputedBlocks[block.Index]; !ok {
			cfg.ComputedBlocks[block.Index] = GetBlockSummary(Context, block)
			deferedFunctions := cfg.ComputedBlocks[block.Index].DeferredFunctions
			cfg.ComputedDeferBlocks[block.Index] = cfg.runDefers(Context, deferedFunctions)
		}

		if len(nextBlock.Succs) == 0 {
			if _, ok := cfg.ComputedBlocks[nextBlock.Index]; !ok {
				cfg.ComputedBlocks[nextBlock.Index] = GetBlockSummary(Context, nextBlock)
				deferedFunctions := cfg.ComputedBlocks[nextBlock.Index].DeferredFunctions
				cfg.ComputedDeferBlocks[nextBlock.Index] = cfg.runDefers(Context, deferedFunctions)
			}
			cfg.visitedBlocksStack.Push(nextBlock)
			cfg.CalculatePath()
			cfg.visitedBlocksStack.Pop()
		} else if !cfg.visitedBlocksStack.Contains(nextBlock) {
			cfg.visitedBlocksStack.Push(nextBlock)
			cfg.traverseGraph(Context, nextBlock)
			cfg.visitedBlocksStack.Pop()
		}
	}

}
