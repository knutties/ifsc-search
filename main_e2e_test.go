package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/blevesearch/bleve/v2"
	"github.com/knutties/ifsc-search/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEndToEnd_CSVThroughHTTP exercises the full pipeline: build a Bleve
// index on disk from a CSV fixture, open it via search.OpenIndex, mount the
// HTTP router, and assert each public endpoint returns sensible data.
func TestEndToEnd_CSVThroughHTTP(t *testing.T) {
	csvPath := filepath.Join("cmd", "build-index", "testdata", "sample.csv")
	indexDir := filepath.Join(t.TempDir(), "index")

	require.NoError(t, buildSmallIndexFromCSV(t, csvPath, indexDir))

	s, err := search.OpenIndex(indexDir)
	require.NoError(t, err)
	t.Cleanup(func() { _ = s.Close() })

	srv := httptest.NewServer(newRouter(s, search.Version{Tag: "test"}, "", io.Discard))
	t.Cleanup(srv.Close)

	t.Run("search", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/search?bank=HDFC&q=andheri")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body search.SearchResults
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.GreaterOrEqual(t, body.Total, 1)
		assert.Equal(t, "HDFC0000001", body.Results[0].IFSC)
	})

	t.Run("healthz", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/healthz")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "ok", body["status"])
		assert.Len(t, body, 1)
	})

	t.Run("status", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/status")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]interface{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "ok", body["status"])
		assert.Equal(t, float64(5), body["indexed_docs"])
		assert.Equal(t, "test", body["release_tag"])
	})

	t.Run("banks", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/list")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Total int           `json:"total"`
			Banks []search.Bank `json:"banks"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, 3, body.Total)
		require.Len(t, body.Banks, 3)
		assert.Equal(t, "HDFC", body.Banks[0].BankCode)
		assert.Equal(t, "HDFC Bank", body.Banks[0].BankName)
		assert.Equal(t, "ICIC", body.Banks[1].BankCode)
		assert.Equal(t, "SBIN", body.Banks[2].BankCode)
	})

	t.Run("lookup_found", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/ifsc/HDFC0000001")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body search.Branch
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "HDFC0000001", body.IFSC)
		assert.Equal(t, "ANDHERI WEST", body.Branch)
		assert.Equal(t, "HDFC Bank", body.BankName)
	})

	t.Run("lookup_not_found", func(t *testing.T) {
		resp, err := http.Get(srv.URL + "/ifsc/ZZZZ0000000")
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

// buildSmallIndexFromCSV mirrors what cmd/build-index does, but kept inline
// here so this test does not depend on importing main from another package.
func buildSmallIndexFromCSV(t *testing.T, csvPath, indexDir string) error {
	t.Helper()
	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	header, err := r.Read()
	if err != nil {
		return err
	}
	cols, err := search.NewColumnIndex(header)
	if err != nil {
		return err
	}

	idx, err := bleve.New(indexDir, search.NewIndexMapping())
	if err != nil {
		return err
	}
	defer idx.Close()

	for {
		row, err := r.Read()
		if err != nil {
			break
		}
		b, err := search.BranchFromCSVRow(cols, row)
		if err != nil {
			continue
		}
		if err := idx.Index(b.IFSC, b); err != nil {
			return err
		}
	}
	return nil
}
