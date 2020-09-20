package main

import (
	"StaticRaceDetector/domain"
	"StaticRaceDetector/utils"
	"fmt"
	"go/token"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
)

func Analysis(pkg *ssa.Package, prog *ssa.Program, accesses []*domain.GuardedAccess) {
	config := &pointer.Config{
		Mains: []*ssa.Package{pkg},
	}

	positionsToGuardAccesses := map[token.Pos][]*domain.GuardedAccess{}
	valuesQueriesToGuardAccess := map[token.Pos]*domain.GuardedAccess{}
	for _, guardedAccess := range accesses {
		if pointer.CanPoint(guardedAccess.Value.Type()) {
			config.AddQuery(guardedAccess.Value)
			valuesQueriesToGuardAccess[guardedAccess.Value.Pos()] = guardedAccess
			positionsToGuardAccesses[guardedAccess.Value.Pos()] = append(positionsToGuardAccesses[guardedAccess.Value.Pos()], guardedAccess)
		}
	}

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

	foundDataRaces := make(map[token.Pos]map[token.Pos]struct{}, 0) // All the data race where the key was already found to avoid duplicates
	for _, guardedAccesses := range positionsToGuardAccesses {
		for _, guardedAccessesA := range guardedAccesses {
			for _, guardedAccessesB := range guardedAccesses {
				if !guardedAccessesA.Intersects(guardedAccessesB) && guardedAccessesA.State.MayConcurrent(guardedAccessesB.State) {
					isExist := utils.DoubleKeyIsExist(guardedAccessesA.Pos, guardedAccessesB.Pos, foundDataRaces) || utils.DoubleKeyIsExist(guardedAccessesB.Pos, guardedAccessesA.Pos, foundDataRaces)
					if !isExist { // If item doesn't exist
						if _ , ok := foundDataRaces[guardedAccessesB.Pos]; !ok {
							foundDataRaces[guardedAccessesB.Pos] = make(map[token.Pos]struct{}, 0)
						}
						foundDataRaces[guardedAccessesB.Pos][guardedAccessesA.Pos] = struct{}{}

						label := fmt.Sprintf(" %s with pos:%s has race condition with %s pos:%s \n", guardedAccessesA.Value, prog.Fset.Position(guardedAccessesA.Pos), guardedAccessesB.Value, prog.Fset.Position(guardedAccessesB.Pos))
						print(label)
					}
				}
			}
		}
	}
}
