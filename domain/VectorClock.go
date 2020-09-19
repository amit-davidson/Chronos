package domain

import (
	"math"
)

type VectorClock map[int]int

func (vc VectorClock) MergeClocks(clockToMerge VectorClock) {
	for goroutine, existingTimestamp := range clockToMerge {
		newTimestamp, _ := vc[goroutine] // 0 as the default value, or the actual value
		vc[goroutine] = int(math.Max(float64(existingTimestamp), float64(newTimestamp)))
	}
}


func (vc VectorClock) Copy() VectorClock {
	copiedClock := make(map[int]int, 0)
	for goroutine, timestamp := range vc {
		copiedClock[goroutine] = timestamp
	}
	return copiedClock
}