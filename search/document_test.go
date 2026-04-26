package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBranchFromCSVRow_PopulatesAllFields(t *testing.T) {
	header := []string{"BANK", "IFSC", "BRANCH", "CENTRE", "DISTRICT", "STATE",
		"ADDRESS", "CONTACT", "IMPS", "RTGS", "CITY", "ISO3166", "NEFT", "MICR",
		"UPI", "SWIFT"}
	row := []string{"HDFC Bank", "HDFC0000240", "ANDHERI WEST", "MUMBAI",
		"MUMBAI", "MAHARASHTRA", "SHOP NO 1-4", "+912226732100", "true", "true",
		"MUMBAI", "IN-MH", "true", "400240015", "true", ""}

	cols, err := NewColumnIndex(header)
	require.NoError(t, err)

	b, err := BranchFromCSVRow(cols, row)
	require.NoError(t, err)

	assert.Equal(t, "HDFC0000240", b.IFSC)
	assert.Equal(t, "HDFC", b.BankCode)
	assert.Equal(t, "HDFC Bank", b.BankName)
	assert.Equal(t, "ANDHERI WEST", b.Branch)
	assert.Equal(t, "MUMBAI", b.Centre)
	assert.Equal(t, "MUMBAI", b.District)
	assert.Equal(t, "MAHARASHTRA", b.State)
	assert.Equal(t, "SHOP NO 1-4", b.Address)
	assert.Equal(t, "MUMBAI", b.City)
	assert.Equal(t, "+912226732100", b.Contact)
	assert.Equal(t, "400240015", b.MICR)
	assert.Equal(t, "", b.SWIFT)
	assert.True(t, b.UPI)
	assert.True(t, b.NEFT)
	assert.True(t, b.RTGS)
	assert.True(t, b.IMPS)
}

func TestBranchFromCSVRow_RejectsShortIFSC(t *testing.T) {
	header := []string{"BANK", "IFSC", "BRANCH", "CENTRE", "DISTRICT", "STATE",
		"ADDRESS", "CONTACT", "IMPS", "RTGS", "CITY", "ISO3166", "NEFT", "MICR",
		"UPI", "SWIFT"}
	row := []string{"X", "BAD", "X", "X", "X", "X", "X", "X", "false", "false",
		"X", "X", "false", "X", "false", ""}
	cols, err := NewColumnIndex(header)
	require.NoError(t, err)

	_, err = BranchFromCSVRow(cols, row)
	assert.Error(t, err, "IFSC shorter than 4 chars must error")
}

func TestNewColumnIndex_RejectsMissingRequiredColumns(t *testing.T) {
	_, err := NewColumnIndex([]string{"BANK", "IFSC"})
	assert.Error(t, err)
}
