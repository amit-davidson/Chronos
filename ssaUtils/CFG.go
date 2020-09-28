package ssaUtils

import (
	"StaticRaceDetector/domain"
	"golang.org/x/tools/go/ssa"
)

type CFG struct {
	ComputedBlockIDsToSummaries map[int]*domain.FunctionState
	visitedBlocks               map[int]struct{}
	lastBlock                   *ssa.BasicBlock
}

func newCFG() *CFG {
	return &CFG{
		ComputedBlockIDsToSummaries: make(map[int]*domain.FunctionState, 0),
		visitedBlocks:               make(map[int]struct{}, 0),
	}
}

func (cfg *CFG) getBlocksSummaries(goroutineState *domain.GoroutineState, entryBlock *ssa.BasicBlock) *domain.FunctionState {
	cfg.getBlocksSummariesDFS(goroutineState, entryBlock)
	funcState := cfg.ComputedBlockIDsToSummaries[cfg.lastBlock.Index]
	return funcState
}

func (cfg *CFG) AreAllPredsCalculated(excludedBackEdges []*ssa.BasicBlock) bool {
	for _, block := range excludedBackEdges {
		_, isExist := cfg.ComputedBlockIDsToSummaries[block.Index]
		if !isExist {
			return false
		}
	}
	return true
}

func (cfg *CFG) mergePredBlocks(blocks []*ssa.BasicBlock) *domain.FunctionState {
	state := cfg.ComputedBlockIDsToSummaries[blocks[0].Index].Copy()
	for _, predBlock := range blocks[1:] {
		predBlockSummary := cfg.ComputedBlockIDsToSummaries[predBlock.Index].Copy()
		state.MergeBlockStates(predBlockSummary)
	}
	return state
}

func (cfg *CFG) getBlocksSummariesDFS(goroutineState *domain.GoroutineState, block *ssa.BasicBlock) {
	// Stop conditions
	_, wasCalculated := cfg.ComputedBlockIDsToSummaries[block.Index]
	if wasCalculated {
		return
	}
	if !cfg.AreAllPredsCalculated(block.Preds) { // If one of the preds wasn't calculated yet, we return and we'll reach this block again once the last pred is calculated. Only then we can merge
		return
	}

	// calculate the merged lockset of all the direct predecessors of the block(Could also be 1). Then merge the old
	// lockset to the new state by updating all the guarded accesses and the lockset at exit point and remove duplicates from the different branches
	var calculatedState *domain.FunctionState
	currBlockState := GetBlockSummary(goroutineState, block)
	if len(block.Preds) > 0 {
		calculatedState = cfg.mergePredBlocks(block.Preds)
		currBlockState.UpdateGuardedAccessesWithLockset(calculatedState.Lockset) // Update all the guarded accesses with the aggregated lockset
		calculatedState.MergeStates(currBlockState)
	} else {
		calculatedState = currBlockState
	}

	cfg.ComputedBlockIDsToSummaries[block.Index] = calculatedState

	if len(block.Succs) == 0 {
		cfg.lastBlock = block
		return
	}

	for _, blockToExecute := range block.Succs {
		cfg.getBlocksSummariesDFS(goroutineState, blockToExecute)
	}
}
