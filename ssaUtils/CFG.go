package ssaUtils

import (
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/utils/stacks"
	"golang.org/x/tools/go/ssa"
)

type CFG struct {
	visitedBlocksStack *stacks.BasicBlockStack

	ComputedBlocks      map[int]*domain.BlockState
	ComputedDeferBlocks map[int]*domain.BlockState
	calculatedState     *domain.BlockState
}

func newCFG() *CFG {
	return &CFG{
		visitedBlocksStack:  stacks.NewBasicBlockStack(),
		ComputedBlocks:      make(map[int]*domain.BlockState),
		ComputedDeferBlocks: make(map[int]*domain.BlockState),
	}
}

func (cfg *CFG) calculateFunctionStatePathSensitive(context *domain.Context, block *ssa.BasicBlock) {
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
	path := cfg.visitedBlocksStack.GetItems()
	block := path[0]
	state := cfg.ComputedBlocks[block.Index].Copy()
	for _, nextBlock := range path[1:] {
		nextState := cfg.ComputedBlocks[nextBlock.Index].Copy()
		state.MergeChildBlock(nextState)
	}

	var firstDeferState *domain.BlockState
	var firstDeferIndex int
	for i := len(path) - 1; i >= 0; i-- {
		deferIndex := path[i].Index
		deferState, ok := cfg.ComputedDeferBlocks[deferIndex]
		if ok {
			firstDeferState = deferState
			firstDeferIndex = i
			break
		}
	}
	if firstDeferState != nil {
		deferStateCopy := firstDeferState.Copy()
		for i := firstDeferIndex - 1; i >= 0; i-- {
			nextState, ok := cfg.ComputedDeferBlocks[path[i].Index]
			if !ok {
				continue
			}
			nextStateCopy := nextState.Copy()
			deferStateCopy.MergeChildBlock(nextStateCopy)
		}
		state.AddResult(deferStateCopy, true)
	}

	if cfg.calculatedState == nil {
		cfg.calculatedState = state
	} else {
		cfg.calculatedState.MergeSiblingBlock(state)
	}
}

func (cfg *CFG) calculateBlockStateIfNeeded(context *domain.Context, block *ssa.BasicBlock) {
	if _, ok := cfg.ComputedBlocks[block.Index]; !ok {
		cfg.ComputedBlocks[block.Index] = GetBlockSummary(context, block)
		deferedFunctions := cfg.ComputedBlocks[block.Index].DeferredFunctions
		if deferedFunctions.Len() > 0 {
			cfg.ComputedDeferBlocks[block.Index] = cfg.runDefers(context, deferedFunctions)
		}
	}
}
