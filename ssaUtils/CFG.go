package ssaUtils

import (
	"StaticRaceDetector/domain"
	"golang.org/x/tools/go/ssa"
)

type CFG struct {
	ComputedBlockIDsToSummaries map[int]*domain.FunctionState
	lastBlock                   *ssa.BasicBlock

	DeferredFunctions                   map[int][]*domain.DeferFunction
	ComputedDeferredBlockIDsToSummaries map[int]*domain.FunctionState
	entryBlock                          *ssa.BasicBlock
}

func newCFG() *CFG {
	return &CFG{
		ComputedBlockIDsToSummaries:         make(map[int]*domain.FunctionState, 0),
		ComputedDeferredBlockIDsToSummaries: make(map[int]*domain.FunctionState, 0),
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

func (cfg *CFG) getDefersSummaries(goroutineState *domain.GoroutineState, entryBlock *ssa.BasicBlock) *domain.FunctionState {
	cfg.calculateDefers(goroutineState, entryBlock)
	funcState := cfg.ComputedDeferredBlockIDsToSummaries[cfg.entryBlock.Index]
	return funcState
}

func (cfg *CFG) runDefers(goroutineState *domain.GoroutineState, defers []*domain.DeferFunction) *domain.FunctionState {
	calculatedState := domain.GetEmptyFunctionState()
	for i := len(defers) - 1; i >= 0; i-- {
		retState := HandleCallCommon(goroutineState, defers[i].Function)
		retState.UpdateGuardedAccessesWithLockset(calculatedState.Lockset)
		calculatedState.MergeStates(retState)
	}
	return calculatedState

}

func (cfg *CFG) AreAllSuccssCalculated(excludedBackEdges []*ssa.BasicBlock) bool {
	for _, block := range excludedBackEdges {
		_, isExist := cfg.ComputedDeferredBlockIDsToSummaries[block.Index]
		if !isExist {
			return false
		}
	}
	return true
}

func (cfg *CFG) mergeSuccssBlocks(blocks []*ssa.BasicBlock) *domain.FunctionState {
	blocksLen := len(blocks)
	state := cfg.ComputedDeferredBlockIDsToSummaries[blocks[blocksLen-1].Index].Copy()
	for i := len(blocks) - 2; i >= 0; i-- {
		predBlockSummary := cfg.ComputedDeferredBlockIDsToSummaries[blocks[i].Index].Copy()
		state.MergeBlockStates(predBlockSummary)
	}
	return state
}

func (cfg *CFG) calculateDefers(goroutineState *domain.GoroutineState, block *ssa.BasicBlock) {
	// Stop conditions
	_, wasCalculated := cfg.ComputedDeferredBlockIDsToSummaries[block.Index]
	if wasCalculated {
		return
	}
	if !cfg.AreAllSuccssCalculated(block.Succs) { // If one of the success wasn't calculated yet, we return and we'll reach this block again once the last pred is calculated. Only then we can merge
		return
	}

	var calculatedState *domain.FunctionState
	currBlockState := cfg.runDefers(goroutineState, cfg.DeferredFunctions[block.Index])
	if len(block.Succs) > 0 {
		calculatedState = cfg.mergeSuccssBlocks(block.Succs)
		currBlockState.UpdateGuardedAccessesWithLockset(calculatedState.Lockset) // Update all the guarded accesses with the aggregated lockset
		calculatedState.MergeStates(currBlockState)
	} else {
		calculatedState = currBlockState
	}

	cfg.ComputedDeferredBlockIDsToSummaries[block.Index] = calculatedState

	if len(block.Preds) == 0 {
		cfg.entryBlock = block
		return
	}

	for _, blockToExecute := range block.Preds {
		cfg.calculateDefers(goroutineState, blockToExecute)
	}

}
