package search

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fixtureBranches returns a small, hand-curated dataset used across tests.
func fixtureBranches() []*Branch {
	return []*Branch{
		{IFSC: "HDFC0000001", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "ANDHERI WEST", City: "MUMBAI", Address: "S V ROAD",
			State: "MAHARASHTRA"},
		{IFSC: "HDFC0000002", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "ANDHERI EAST", City: "MUMBAI", Address: "CHAKALA",
			State: "MAHARASHTRA"},
		{IFSC: "HDFC0000003", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "BANDRA", City: "MUMBAI", Address: "LINKING ROAD",
			State: "MAHARASHTRA"},
		{IFSC: "ICIC0000001", BankCode: "ICIC", BankName: "ICICI Bank",
			Branch: "ANDHERI WEST", City: "MUMBAI", Address: "JUHU LANE",
			State: "MAHARASHTRA"},
		{IFSC: "SBIN0000001", BankCode: "SBIN", BankName: "State Bank of India",
			Branch: "KOREGAON PARK", City: "PUNE", Address: "NORTH MAIN ROAD",
			State: "MAHARASHTRA"},
		{IFSC: "HDFC0000004", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "ANDHERI NORTH", City: "MUMBAI", Address: "ANDHERI MAIN ROAD",
			State: "MAHARASHTRA"},
	}
}

func newTestSearcher(t *testing.T) Searcher {
	t.Helper()
	s, err := NewMemorySearcher(fixtureBranches())
	require.NoError(t, err)
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestSearchRequest_Validate_RejectsBothEmpty(t *testing.T) {
	req := SearchRequest{}
	err := req.Validate()
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrMissingQuery), "want ErrMissingQuery, got %v", err)
}

func TestSearchRequest_Validate_RejectsNegativeOffset(t *testing.T) {
	req := SearchRequest{Q: "andheri", Offset: -1}
	err := req.Validate()
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrBadPagination))
}

func TestSearchRequest_NormalizePagination_ClampsLimit(t *testing.T) {
	req := SearchRequest{Q: "x", Limit: 5000}
	req.normalize()
	assert.Equal(t, 100, req.Limit, "limit clamped to max")

	req2 := SearchRequest{Q: "x", Limit: 0}
	req2.normalize()
	assert.Equal(t, 20, req2.Limit, "zero falls back to default")
}

func TestSearch_BankFilter_ExactCode(t *testing.T) {
	s := newTestSearcher(t)
	res, err := s.Search(SearchRequest{Bank: "HDFC", Q: "andheri"})
	require.NoError(t, err)
	require.GreaterOrEqual(t, res.Total, 2)
	for _, r := range res.Results {
		assert.Equal(t, "HDFC", r.BankCode)
	}
}

func TestSearch_BankFilter_FuzzyName(t *testing.T) {
	s := newTestSearcher(t)
	res, err := s.Search(SearchRequest{Bank: "ICICI", Q: "andheri"})
	require.NoError(t, err)
	require.Len(t, res.Results, 1)
	assert.Equal(t, "ICIC0000001", res.Results[0].IFSC)
}

func TestSearch_FreeText_NoBank_RanksBranchOverAddress(t *testing.T) {
	s := newTestSearcher(t)
	res, err := s.Search(SearchRequest{Q: "andheri"})
	require.NoError(t, err)
	require.GreaterOrEqual(t, res.Total, 3)
	// The three "ANDHERI" branches should appear before any address-only match.
	assert.Contains(t, []string{"HDFC0000001", "HDFC0000002", "ICIC0000001", "HDFC0000004"},
		res.Results[0].IFSC)
}

func TestSearch_FuzzyTypo(t *testing.T) {
	s := newTestSearcher(t)
	res, err := s.Search(SearchRequest{Q: "andehri"})
	require.NoError(t, err)
	assert.Greater(t, res.Total, 0, "fuzzy match should still find ANDHERI branches")
}

func TestSearch_BankOnly_ReturnsAllBranchesForBank(t *testing.T) {
	s := newTestSearcher(t)
	res, err := s.Search(SearchRequest{Bank: "HDFC"})
	require.NoError(t, err)
	assert.Equal(t, 4, res.Total)
}

func TestSearch_NoMatch_ReturnsEmpty(t *testing.T) {
	s := newTestSearcher(t)
	res, err := s.Search(SearchRequest{Bank: "HDFC", Q: "xyzzyqqq"})
	require.NoError(t, err)
	assert.Equal(t, 0, res.Total)
	assert.Empty(t, res.Results)
}

func TestSearch_UnknownBank_ReturnsEmpty(t *testing.T) {
	s := newTestSearcher(t)
	res, err := s.Search(SearchRequest{Bank: "ZZZZ", Q: "andheri"})
	require.NoError(t, err)
	assert.Equal(t, 0, res.Total)
	assert.Empty(t, res.Results)
}

func TestSearch_PaginationOffsetAndTotal(t *testing.T) {
	s := newTestSearcher(t)
	res, err := s.Search(SearchRequest{Q: "andheri", Limit: 1, Offset: 0})
	require.NoError(t, err)
	assert.Len(t, res.Results, 1)
	totalFirstPage := res.Total

	res2, err := s.Search(SearchRequest{Q: "andheri", Limit: 1, Offset: 1})
	require.NoError(t, err)
	assert.Equal(t, totalFirstPage, res2.Total, "total stable across pages")
	if len(res2.Results) > 0 {
		assert.NotEqual(t, res.Results[0].IFSC, res2.Results[0].IFSC)
	}
}
