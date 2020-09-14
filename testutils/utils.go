package testutils

import (
	"StaticRaceDetector/domain"
	"StaticRaceDetector/utils"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestResult struct {
	Lockset *domain.Lockset
	GuardedAccess []*domain.GuardedAccess
}

func WriteResult(t *testing.T, path string, ls *domain.Lockset, ga []*domain.GuardedAccess) {
	testresult := TestResult{Lockset: ls, GuardedAccess: ga}
	dump, err := json.Marshal(testresult)
	require.NoError(t, err)
	utils.UpdateFile(t, path, dump)
}
