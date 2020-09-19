package main

import (
	"StaticRaceDetector/domain"
	"fmt"
	"go/token"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
)

func Analysis(pkg *ssa.Package, prog *ssa.Program, accesses []*domain.GuardedAccess) {
	config := &pointer.Config{
		Mains: []*ssa.Package{pkg},
	}

	valuesQueriesToGuardAccess := map[token.Pos]*domain.GuardedAccess{}
	for _, guardedAccess := range accesses {
		if pointer.CanPoint(guardedAccess.Value.Type()) {
			config.AddQuery(guardedAccess.Value)
			valuesQueriesToGuardAccess[guardedAccess.Value.Pos()] = guardedAccess
		}
	}

	positionsToGuardAccesses := map[token.Pos][]*domain.GuardedAccess{}
	result, err := pointer.Analyze(config)
	if err != nil {
		panic(err) // internal error in pointer analysis
	}

	for v, l := range result.Queries {
		for _, label := range l.PointsTo().Labels() {
			guardedAccess := valuesQueriesToGuardAccess[v.Pos()]
			allocPos := label.Value()
			positionsToGuardAccesses[allocPos.Pos()] = append(positionsToGuardAccesses[allocPos.Pos()], guardedAccess)
		}
	}
	for _, guardedAccesses := range positionsToGuardAccesses {
		for _, guardedAccessesA := range guardedAccesses {
			for _, guardedAccessesB := range guardedAccesses {
				if !guardedAccessesA.Intersects(guardedAccessesB) && guardedAccessesA.State.MayConcurrent(guardedAccessesB.State) {
					valueA := guardedAccessesA.Value
					valueB := guardedAccessesB.Value
					label := fmt.Sprintf(" %s with pos:%s has race condition with %s pos:%s \n", valueA, prog.Fset.Position(valueA.Pos()), valueB, prog.Fset.Position(valueB.Pos()))
					print(label)
				}
			}
		}
	}

}
