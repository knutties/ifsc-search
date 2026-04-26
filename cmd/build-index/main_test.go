package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/razorpay/ifsc/v2/ifsc-api/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildIndexFromCSV_RoundTripsAllRows(t *testing.T) {
	csvPath := filepath.Join("testdata", "sample.csv")
	indexDir := filepath.Join(t.TempDir(), "index")

	err := buildIndexFromCSV(csvPath, indexDir)
	require.NoError(t, err)

	s, err := search.OpenIndex(indexDir)
	require.NoError(t, err)
	defer s.Close()
	assert.Equal(t, uint64(5), s.DocCount())

	res, err := s.Search(search.SearchRequest{Bank: "HDFC"})
	require.NoError(t, err)
	assert.Equal(t, 3, res.Total, "three HDFC rows in the fixture")
}

func TestBuildIndexFromCSV_RejectsBadHeader(t *testing.T) {
	csvPath := filepath.Join(t.TempDir(), "bad.csv")
	require.NoError(t, os.WriteFile(csvPath,
		[]byte("FOO,BAR\n1,2\n"), 0644))
	err := buildIndexFromCSV(csvPath, filepath.Join(t.TempDir(), "idx"))
	assert.Error(t, err)
}
