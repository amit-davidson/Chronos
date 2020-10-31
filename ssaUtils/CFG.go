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

func (cfg *CFG) calculateFunctionState(context *domain.Context, block *ssa.BasicBlock) {
	firstBlock := &ssa.BasicBlock{Index: -1, Succs: []*ssa.BasicBlock{block}}
	cfg.traverseGraph(context, firstBlock)
}

func (cfg *CFG) traverseGraph(context *domain.Context, block *ssa.BasicBlock) {
	nextBlocks := block.Succs
	for _, nextBlock := range nextBlocks {
		cfg.calculateBlockStateIfNeeded(context, block)
		if len(nextBlock.Succs) == 0 {
			cfg.calculateBlockStateIfNeeded(context, nextBlock)
			cfg.visitedBlocksStack.Push(nextBlock)
			cfg.CalculatePath()
			cfg.visitedBlocksStack.Pop()
		} else if !cfg.visitedBlocksStack.Contains(nextBlock) {
			cfg.visitedBlocksStack.Push(nextBlock)
			cfg.traverseGraph(context, nextBlock)
			cfg.visitedBlocksStack.Pop()
		}
	}
}

func (cfg *CFG) CalculatePath() {
	path := cfg.visitedBlocksStack.GetAllItems()
	block := path[0]
	state := cfg.ComputedBlocks[block.Index].Copy()
	for _, nextBlock := range path[1:] {
		nextState := cfg.ComputedBlocks[nextBlock.Index].Copy()
		nextState.UpdateGuardedAccessesWithLockset(state.Lockset)
		state.MergeStates(nextState, true)
	}

	deferBlock := path[len(path)-1]
	deferState, ok := cfg.ComputedDeferBlocks[deferBlock.Index]
	if ok {
		deferStateCopy := deferState.Copy()
		for i := len(path) - 2; i >= 0; i-- {
			nextState, ok := cfg.ComputedDeferBlocks[path[i].Index]
			if !ok {
				continue
			}
			nextStateCopy := nextState.Copy()
			nextStateCopy.UpdateGuardedAccessesWithLockset(deferStateCopy.Lockset)
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

func (cfg *CFG) calculateBlockStateIfNeeded(context *domain.Context, block *ssa.BasicBlock) {
	if _, ok := cfg.ComputedBlocks[block.Index]; !ok {
		cfg.ComputedBlocks[block.Index] = GetBlockSummary(context, block)
		deferedFunctions := cfg.ComputedBlocks[block.Index].DeferredFunctions
		cfg.ComputedDeferBlocks[block.Index] = cfg.runDefers(context, deferedFunctions)
	}
}
