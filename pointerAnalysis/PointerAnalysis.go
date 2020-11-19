package pointerAnalysis

import (
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/utils"
	"go/token"
	"golang.org/x/tools/go/pointer"
	"golang.org/x/tools/go/ssa"
)


// Analysis starts by mapping between positions of the guard accesses (values inside) to the guard accesses themselves.
// Then it analyzes all the values inside values inside, and check if some of the values might alias each other. If so,
// the positions for those values are merged. After all positions were merged, the algorithm runs and check if for a
// given value (identified by a pos in the map) there are two guarded accesses that might conflict - W&W/R&W from two
// different goroutines.
// map1 : A->ga1, ga2, ga3
//        B->ga4, ga5, ga6
//        C->ga7, ga8, ga9
//        D->ga10, ga11, ga12
// Now that we know that B may point to A, we add it to it
// map1 : A->ga1, ga2, ga3, ga4, ga5, ga6
//        C->ga7, ga8, ga9
//        D->ga10, ga11, ga12
// And if A may point to D, then
// map1 : C->ga7, ga8, ga9
//        D->ga10, ga11, ga12, ga1, ga2, ga3, ga4, ga5, ga6
// And then for pos all the guarded accesses are compared to see if data races might exist

func Analysis(pkg *ssa.Package, accesses []*domain.GuardedAccess) ([][]*domain.GuardedAccess, error) {
	config := &pointer.Config{
		Mains: []*ssa.Package{pkg},
	}

	positionsToGuardAccesses := map[token.Pos][]*domain.GuardedAccess{}
	for _, guardedAccess := range accesses {
		if guardedAccess.Pos.IsValid() && pointer.CanPoint(guardedAccess.Value.Type()) {
			config.AddQuery(guardedAccess.Value)
			// Multiple instructions for the same variable for example write and multiple reads
			positionsToGuardAccesses[guardedAccess.Value.Pos()] = append(positionsToGuardAccesses[guardedAccess.Value.Pos()], guardedAccess)
		}
	}

	result, err := pointer.Analyze(config)
	if err != nil {
		return nil, err // internal error in pointer analysis
	}

	// Join instructions of variables that may point to each other.
	for v, l := range result.Queries {
		for _, label := range l.PointsTo().Labels() {
			allocPos := label.Value().Pos()
			queryPos := v.Pos()
			if allocPos == queryPos {
				continue
			}
			positionsToGuardAccesses[allocPos] = append(positionsToGuardAccesses[allocPos], positionsToGuardAccesses[queryPos]...)
		}
	}
	foundDataRaces := utils.NewDoubleKeyMap() // To avoid reporting on the same pair of positions more then once. Can happen if for the same place we read and then write.
	conflictingGA := make([][]*domain.GuardedAccess, 0)
	for _, guardedAccesses := range positionsToGuardAccesses {
		for _, guardedAccessA := range guardedAccesses {
			for _, guardedAccessB := range guardedAccesses {
				if guardedAccessA.IsConflicting(guardedAccessB) {
					isExist := foundDataRaces.IsExist(guardedAccessA.Pos, guardedAccessB.Pos)
					if !isExist {
						foundDataRaces.Add(guardedAccessA.Pos, guardedAccessB.Pos)
						conflictingGA = append(conflictingGA, []*domain.GuardedAccess{guardedAccessA, guardedAccessB})
					}
				}
			}
		}
	}
	return conflictingGA, nil
}

