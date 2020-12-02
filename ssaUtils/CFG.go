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
	FunctionState       *domain.BlockState
}

func newCFG() *CFG {
	return &CFG{
		visitedBlocksStack:  stacks.NewBlockMap(),
		ComputedBlocks:      make(map[int]*domain.BlockState),
		ComputedDeferBlocks: make(map[int]*domain.BlockState),
		FunctionState:       domain.GetEmptyBlockState(),
	}
}

// CalculateFunctionState works by traversing the tree in a DFS way, similar to the flow of the function when it'll run.
// It calculates the state of the regular flow of block first, then adds the state of any succeeding blocks in the tree,
// and finally the block's defer state if exist.
// The function uses two way to aggregate the states between blocks. If the blocks are adjacent (siblings) to each
// other, (resulted from a branch) then a merge mechanism is used. If one block is below the other, then an append is
// performed.
func (cfg *CFG) CalculateFunctionState(context *domain.Context, block *ssa.BasicBlock) {
	cfg.visitedBlocksStack.Add(block)
	defer cfg.visitedBlocksStack.Remove(block)
	cfg.calculateBlockState(context, block)

	// Regular flow
	blockState := cfg.ComputedBlocks[block.Index]
	cfg.FunctionState.MergeSiblingBlock(blockState)

	// recursion
	for _, nextBlock := range block.Succs {
		// if it's a cycle we skip it
		if cfg.visitedBlocksStack.Contains(nextBlock.Index) {
			continue
		}

		cfg.CalculateFunctionState(context, nextBlock)
	}

	// Defer
	if deferState, ok := cfg.ComputedDeferBlocks[block.Index]; ok {
		blockState.MergeChildBlock(deferState)
	}
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
