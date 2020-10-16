package testutils

import (
	"encoding/json"
	"github.com/amit-davidson/Chronos/domain"
	"github.com/amit-davidson/Chronos/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestResult struct {
	Lockset       *LocksetJson
	GuardedAccess []GuardedAccessJSON
}

func NewTestResult(ls *domain.Lockset, gal []*domain.GuardedAccess) *TestResult {
	ts := &TestResult{Lockset: LocksetToJSON(ls)}
	guardedAccessesJson := make([]GuardedAccessJSON, 0)
	for _, ga := range gal {
		guardedAccessesJson = append(guardedAccessesJson, GuardedAccessToJSON(ga))
	}
	ts.GuardedAccess = guardedAccessesJson
	return ts
}

func WriteResult(t *testing.T, path string, ls *domain.Lockset, ga []*domain.GuardedAccess) {
	testresult := NewTestResult(ls, ga)
	dump, err := json.Marshal(testresult)
	require.NoError(t, err)
	utils.UpdateFile(t, path, dump)
}

func CompareResult(t *testing.T, path string, ls *domain.Lockset, ga []*domain.GuardedAccess) {
	testresult := &TestResult{}
	data, err := utils.ReadFile(path)
	require.NoError(t, err)
	err = json.Unmarshal(data, testresult)
	require.NoError(t, err)

	assert.Equal(t, testresult.Lockset, LocksetToJSON(ls))
	assert.Equal(t, len(testresult.GuardedAccess), len(ga))
	for i := range ga {
		insr := GuardedAccessToJSON(ga[i])
		assert.Equal(t, testresult.GuardedAccess[i], insr)
	}
}
