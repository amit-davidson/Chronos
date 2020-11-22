package ssaUtils

import (
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/utils/stacks"
	"golang.org/x/tools/go/ssa"
)

type CFG struct {
	visitedBlocksStack *stacks.BlockMap

	ComputedBlocks      map[int]*domain.BlockState
	ComputedDeferBlocks map[int]*domain.BlockState
}

func newCFG() *CFG {
	return &CFG{
		visitedBlocksStack:  stacks.NewBlockMap(),
		ComputedBlocks:      make(map[int]*domain.BlockState),
		ComputedDeferBlocks: make(map[int]*domain.BlockState),
	}
}

func (cfg *CFG) CalculateFunctionStatePathSensitive(context *domain.Context, block *ssa.BasicBlock) *domain.BlockState {
	cfg.visitedBlocksStack.Add(block)
	defer cfg.visitedBlocksStack.Remove(block)
	cfg.calculateBlockState(context, block)

	// Regular flow
	blockState := cfg.ComputedBlocks[block.Index]

	// recursion
	var branchState *domain.BlockState
	for _, nextBlock := range block.Succs {
		// if it's a cycle we skip it
		if cfg.visitedBlocksStack.Contains(nextBlock.Index) {
			continue
		}
		retBlockState := cfg.CalculateFunctionStatePathSensitive(context, nextBlock)
		if branchState == nil {
			branchState = retBlockState.Copy()
		} else {
			branchState.MergeSiblingBlock(retBlockState)
		}
	}
	if branchState != nil {
		blockState.MergeChildBlock(branchState)
	}

	// Defer
	if deferState, ok := cfg.ComputedDeferBlocks[block.Index]; ok {
		blockState.MergeChildBlock(deferState)
	}
	return blockState
}

func (cfg *CFG) calculateBlockState(context *domain.Context, block *ssa.BasicBlock) {
	if _, ok := cfg.ComputedBlocks[block.Index]; !ok {
		cfg.ComputedBlocks[block.Index] = GetBlockSummary(context, block)
		deferedFunctions := cfg.ComputedBlocks[block.Index].DeferredFunctions
		if deferedFunctions.Len() > 0 {
			cfg.ComputedDeferBlocks[block.Index] = cfg.runDefers(context, deferedFunctions)
		}
	}
}
