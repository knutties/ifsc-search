package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/knutties/ifsc-search/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	branches := []*search.Branch{
		{IFSC: "HDFC0000001", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "ANDHERI WEST", City: "MUMBAI", State: "MAHARASHTRA"},
		{IFSC: "HDFC0000002", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "BANDRA", City: "MUMBAI", State: "MAHARASHTRA"},
		{IFSC: "ICIC0000001", BankCode: "ICIC", BankName: "ICICI Bank",
			Branch: "ANDHERI EAST", City: "MUMBAI", State: "MAHARASHTRA"},
	}
	s, err := search.NewMemorySearcher(branches)
	require.NoError(t, err)
	t.Cleanup(func() { _ = s.Close() })

	srv := httptest.NewServer(newRouter(s, search.Version{
		Tag: "test", BuiltAt: "2026-04-26T00:00:00Z"}, ""))
	t.Cleanup(srv.Close)
	return srv
}

func TestHandleSearch_BankAndQuery(t *testing.T) {
	srv := newTestServer(t)
	resp, err := http.Get(srv.URL + "/search?bank=HDFC&q=andheri")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body search.SearchResults
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, 1, body.Total)
	assert.Equal(t, "HDFC0000001", body.Results[0].IFSC)
}

func TestHandleSearch_MissingParams_400(t *testing.T) {
	srv := newTestServer(t)
	resp, err := http.Get(srv.URL + "/search")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var body map[string]string
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Contains(t, body["error"], "at least one of bank or q is required")
}

func TestHandleSearch_NegativeOffset_400(t *testing.T) {
	srv := newTestServer(t)
	resp, err := http.Get(srv.URL + "/search?q=andheri&offset=-1")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHandleSearch_NonIntegerLimit_400(t *testing.T) {
	srv := newTestServer(t)
	resp, err := http.Get(srv.URL + "/search?q=andheri&limit=abc")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHandleHealthz(t *testing.T) {
	srv := newTestServer(t)
	resp, err := http.Get(srv.URL + "/healthz")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]interface{}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "ok", body["status"])
	assert.Equal(t, float64(3), body["indexed_docs"])
	assert.Equal(t, "test", body["release_tag"])
}

func TestHandleListBanks(t *testing.T) {
	srv := newTestServer(t)
	resp, err := http.Get(srv.URL + "/banks")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body struct {
		Total int           `json:"total"`
		Banks []search.Bank `json:"banks"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, 2, body.Total)
	require.Len(t, body.Banks, 2)
	assert.Equal(t, "HDFC", body.Banks[0].BankCode)
	assert.Equal(t, "HDFC Bank", body.Banks[0].BankName)
	assert.Equal(t, "ICIC", body.Banks[1].BankCode)
	assert.Equal(t, "ICICI Bank", body.Banks[1].BankName)
}

func TestHandleLookup_Found(t *testing.T) {
	srv := newTestServer(t)
	resp, err := http.Get(srv.URL + "/ifsc/HDFC0000001")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body search.Branch
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "HDFC0000001", body.IFSC)
	assert.Equal(t, "ANDHERI WEST", body.Branch)
}

func TestHandleLookup_NotFound(t *testing.T) {
	srv := newTestServer(t)
	resp, err := http.Get(srv.URL + "/ifsc/ZZZZ0000000")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestHandleLookup_MethodNotAllowed(t *testing.T) {
	srv := newTestServer(t)
	req, err := http.NewRequest(http.MethodPost, srv.URL+"/ifsc/HDFC0000001", nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}

func TestRouter_PathPrefix(t *testing.T) {
	branches := []*search.Branch{
		{IFSC: "HDFC0000001", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "ANDHERI WEST", City: "MUMBAI", State: "MAHARASHTRA"},
	}
	s, err := search.NewMemorySearcher(branches)
	require.NoError(t, err)
	t.Cleanup(func() { _ = s.Close() })

	srv := httptest.NewServer(newRouter(s, search.Version{Tag: "test"}, "/ifsc"))
	t.Cleanup(srv.Close)

	resp, err := http.Get(srv.URL + "/ifsc/search?bank=HDFC&q=andheri")
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = http.Get(srv.URL + "/ifsc/healthz")
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = http.Get(srv.URL + "/ifsc/ifsc/HDFC0000001")
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = http.Get(srv.URL + "/search?bank=HDFC&q=andheri")
	require.NoError(t, err)
	resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestNormalizePrefix(t *testing.T) {
	cases := map[string]string{
		"":         "",
		"/":        "",
		"ifsc":     "/ifsc",
		"/ifsc":    "/ifsc",
		"/ifsc/":   "/ifsc",
		"ifsc/v1":  "/ifsc/v1",
		"/ifsc/v1": "/ifsc/v1",
		"  /api ":  "/api",
	}
	for in, want := range cases {
		assert.Equal(t, want, normalizePrefix(in), "input=%q", in)
	}
}
