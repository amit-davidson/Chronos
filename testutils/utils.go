package testutils

import "StaticRaceDetector/domain"

type testResult struct {
	ls *
	ga []*guardedAccess
}

func WriteResult(ls *domain.Lockset, ga []*domain.GuardedAccess) ([]byte, error) {

}