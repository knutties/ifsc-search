package search

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fixtureBranches returns a small, hand-curated dataset used across tests.
func fixtureBranches() []*Branch {
	return []*Branch{
		{IFSC: "HDFC0000001", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "ANDHERI WEST", City: "MUMBAI", Address: "S V ROAD",
			District: "MUMBAI SUBURBAN", State: "MAHARASHTRA"},
		{IFSC: "HDFC0000002", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "ANDHERI EAST", City: "MUMBAI", Address: "CHAKALA",
			District: "MUMBAI SUBURBAN", State: "MAHARASHTRA"},
		{IFSC: "HDFC0000003", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "BANDRA", City: "MUMBAI", Address: "LINKING ROAD",
			District: "MUMBAI SUBURBAN", State: "MAHARASHTRA"},
		{IFSC: "ICIC0000001", BankCode: "ICIC", BankName: "ICICI Bank",
			Branch: "ANDHERI WEST", City: "MUMBAI", Address: "JUHU LANE",
			District: "MUMBAI SUBURBAN", State: "MAHARASHTRA"},
		{IFSC: "SBIN0000001", BankCode: "SBIN", BankName: "State Bank of India",
			Branch: "KOREGAON PARK", City: "PUNE", Address: "NORTH MAIN ROAD",
			District: "PUNE", State: "MAHARASHTRA"},
		{IFSC: "HDFC0000004", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "ANDHERI NORTH", City: "MUMBAI", Address: "ANDHERI MAIN ROAD",
			District: "MUMBAI SUBURBAN", State: "MAHARASHTRA"},
		// Extra rows that vary state/district/city so strict-equality
		// filters can be exercised without extra setup.
		{IFSC: "SBIN0000002", BankCode: "SBIN", BankName: "State Bank of India",
			Branch: "VASHI", City: "NAVI MUMBAI", Address: "PALM BEACH ROAD",
			District: "THANE", State: "MAHARASHTRA"},
		{IFSC: "SBIN0000003", BankCode: "SBIN", BankName: "State Bank of India",
			Branch: "MG ROAD", City: "BANGALORE", Address: "MG ROAD",
			District: "BANGALORE URBAN", State: "KARNATAKA"},
		{IFSC: "SBIN0000004", BankCode: "SBIN", BankName: "State Bank of India",
			Branch: "PARK STREET", City: "KOLKATA", Address: "PARK STREET",
			District: "KOLKATA", State: "WEST BENGAL"},
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

func TestSearchRequest_Validate_AcceptsAnyOneSignal(t *testing.T) {
	for name, req := range map[string]SearchRequest{
		"ifsc":     {IFSCPrefix: "HDFC0"},
		"state":    {State: "Maharashtra"},
		"district": {District: "Mumbai Suburban"},
		"city":     {City: "Mumbai"},
	} {
		t.Run(name, func(t *testing.T) {
			assert.NoError(t, req.Validate())
		})
	}
}

func TestSearch_IFSCPrefix_Matches(t *testing.T) {
	s := newTestSearcher(t)
	res, err := s.Search(SearchRequest{IFSCPrefix: "HDFC0"})
	require.NoError(t, err)
	assert.Equal(t, 4, res.Total, "all four HDFC branches should match")
	for _, r := range res.Results {
		assert.True(t, strings.HasPrefix(r.IFSC, "HDFC0"), "got %s", r.IFSC)
	}
}

func TestSearch_IFSCPrefix_CaseInsensitive(t *testing.T) {
	s := newTestSearcher(t)
	upper, err := s.Search(SearchRequest{IFSCPrefix: "SBIN0"})
	require.NoError(t, err)
	lower, err := s.Search(SearchRequest{IFSCPrefix: "sbin0"})
	require.NoError(t, err)
	assert.Equal(t, upper.Total, lower.Total)
	assert.Greater(t, upper.Total, 0)
}

func TestSearch_City_StrictDoesNotMatchNaviMumbai(t *testing.T) {
	s := newTestSearcher(t)
	res, err := s.Search(SearchRequest{City: "Mumbai"})
	require.NoError(t, err)
	for _, r := range res.Results {
		assert.Equal(t, "MUMBAI", r.City,
			"city=Mumbai must not bleed into NAVI MUMBAI: got %s", r.City)
	}
	// Five Mumbai-proper branches in the fixture; SBIN0000002 (Navi Mumbai)
	// must be excluded.
	assert.Equal(t, 5, res.Total)
}

func TestSearch_State_SingleWordExact(t *testing.T) {
	s := newTestSearcher(t)
	res, err := s.Search(SearchRequest{State: "Karnataka"})
	require.NoError(t, err)
	require.Equal(t, 1, res.Total)
	assert.Equal(t, "SBIN0000003", res.Results[0].IFSC)
}

func TestSearch_State_MultiWordExact(t *testing.T) {
	s := newTestSearcher(t)
	res, err := s.Search(SearchRequest{State: "West Bengal"})
	require.NoError(t, err)
	require.Equal(t, 1, res.Total, "multi-word state must resolve via keyword field, not tokenized text")
	assert.Equal(t, "SBIN0000004", res.Results[0].IFSC)
}

func TestSearch_District_MultiWord(t *testing.T) {
	s := newTestSearcher(t)
	res, err := s.Search(SearchRequest{District: "Bangalore Urban"})
	require.NoError(t, err)
	require.Equal(t, 1, res.Total)
	assert.Equal(t, "SBIN0000003", res.Results[0].IFSC)
}

func TestSearch_Q_MatchesIFSCCode(t *testing.T) {
	s := newTestSearcher(t)
	res, err := s.Search(SearchRequest{Q: "HDFC0000003"})
	require.NoError(t, err)
	require.GreaterOrEqual(t, res.Total, 1)
	assert.Equal(t, "HDFC0000003", res.Results[0].IFSC,
		"q with a full IFSC must surface that branch first")
}

func TestSearch_Q_MatchesIFSCPrefixCaseInsensitive(t *testing.T) {
	s := newTestSearcher(t)
	res, err := s.Search(SearchRequest{Q: "sbin0"})
	require.NoError(t, err)
	require.GreaterOrEqual(t, res.Total, 4, "all four SBIN branches should surface via q prefix")
	for _, r := range res.Results[:4] {
		assert.True(t, strings.HasPrefix(r.IFSC, "SBIN"), "got %s", r.IFSC)
	}
}

func TestSearch_FilterCombinations_AndSemantics(t *testing.T) {
	s := newTestSearcher(t)
	// HDFC + Mumbai → 4. HDFC + Karnataka → 0 (HDFC has no Karnataka branch).
	hdfcMumbai, err := s.Search(SearchRequest{Bank: "HDFC", City: "Mumbai"})
	require.NoError(t, err)
	assert.Equal(t, 4, hdfcMumbai.Total)

	hdfcKarnataka, err := s.Search(SearchRequest{Bank: "HDFC", State: "Karnataka"})
	require.NoError(t, err)
	assert.Equal(t, 0, hdfcKarnataka.Total)
}
