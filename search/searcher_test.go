package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMemorySearcher_IndexesAndCounts(t *testing.T) {
	branches := []*Branch{
		{IFSC: "HDFC0000001", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "ANDHERI WEST", City: "MUMBAI", State: "MAHARASHTRA"},
		{IFSC: "HDFC0000002", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "BANDRA", City: "MUMBAI", State: "MAHARASHTRA"},
	}
	s, err := NewMemorySearcher(branches)
	require.NoError(t, err)
	defer s.Close()

	assert.Equal(t, uint64(2), s.DocCount())
}
