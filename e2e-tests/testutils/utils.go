package testutils

import (
	"Miranda/domain"
	"Miranda/utils"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type TestResult struct {
	Lockset       *domain.Lockset
	GuardedAccess []*domain.GuardedAccess
}

type TestResultJSON struct {
	Lockset       *domain.LocksetJson
	GuardedAccess []domain.GuardedAccessJSON
}

func WriteResult(t *testing.T, path string, ls *domain.Lockset, ga []*domain.GuardedAccess) {
	testresult := TestResult{Lockset: ls, GuardedAccess: ga}
	dump, err := json.Marshal(testresult)
	require.NoError(t, err)
	utils.UpdateFile(t, path, dump)
}

func CompareResult(t *testing.T, path string, ls *domain.Lockset, ga []*domain.GuardedAccess) {
	testresult := &TestResultJSON{}
	data, err := utils.ReadFile(path)
	require.NoError(t, err)
	err = json.Unmarshal(data, testresult)
	require.NoError(t, err)

	assert.Equal(t, testresult.Lockset, ls.ToJSON())
	assert.Equal(t, len(testresult.GuardedAccess), len(ga))
	for i := range ga {
		insr := ga[i].ToJSON()
		assert.Equal(t, testresult.GuardedAccess[i], insr)
	}
}
