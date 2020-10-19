package domain

import (
	"math"
)

type VectorClock map[int]int

func (vc VectorClock) MergeClocks(clockToMerge VectorClock) {
	for goroutine, existingTimestamp := range clockToMerge {
		newTimestamp := vc.Get(goroutine) // 0 as the default value, or the actual value
		vc[goroutine] = int(math.Max(float64(existingTimestamp), float64(newTimestamp)))
	}
}

// Get returns clock for provided goroutine ID.
// If no such goroutine or clock is not initialized,
// then returns a zero value.
func (vc VectorClock) Get(id int) int {
	if vc == nil {
		return 0
	}
	var ts, _ = vc[id]
	return ts
}

func (vc VectorClock) Copy() VectorClock {
	copiedClock := make(map[int]int, len(vc))
	for goroutine, timestamp := range vc {
		copiedClock[goroutine] = timestamp
	}
	return copiedClock
}
