package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/knutties/ifsc-search/search"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return newServerWith(t, []*search.Branch{
		{IFSC: "HDFC0000001", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "ANDHERI WEST", City: "MUMBAI", State: "MAHARASHTRA"},
		{IFSC: "HDFC0000002", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "BANDRA", City: "MUMBAI", State: "MAHARASHTRA"},
		{IFSC: "ICIC0000001", BankCode: "ICIC", BankName: "ICICI Bank",
			Branch: "ANDHERI EAST", City: "MUMBAI", State: "MAHARASHTRA"},
	})
}

func newServerWith(t *testing.T, branches []*search.Branch) *httptest.Server {
	t.Helper()
	s, err := search.NewMemorySearcher(branches)
	require.NoError(t, err)
	t.Cleanup(func() { _ = s.Close() })

	srv := httptest.NewServer(newRouter(s, search.Version{
		Tag: "test", BuiltAt: "2026-04-26T00:00:00Z"}, "", io.Discard))
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
	assert.Contains(t, body["error"], "at least one of bank, q, ifsc, state, district, city is required")
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

func filterFixture() []*search.Branch {
	return []*search.Branch{
		{IFSC: "HDFC0000001", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "ANDHERI WEST", City: "MUMBAI",
			District: "MUMBAI SUBURBAN", State: "MAHARASHTRA"},
		{IFSC: "SBIN0000001", BankCode: "SBIN", BankName: "State Bank of India",
			Branch: "VASHI", City: "NAVI MUMBAI",
			District: "THANE", State: "MAHARASHTRA"},
		{IFSC: "SBIN0000002", BankCode: "SBIN", BankName: "State Bank of India",
			Branch: "MG ROAD", City: "BANGALORE",
			District: "BANGALORE URBAN", State: "KARNATAKA"},
	}
}

func TestHandleSearch_IFSCPrefix(t *testing.T) {
	srv := newServerWith(t, filterFixture())
	resp, err := http.Get(srv.URL + "/search?ifsc=SBIN0")
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body search.SearchResults
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, 2, body.Total)
	for _, r := range body.Results {
		assert.Equal(t, "SBIN", r.BankCode)
	}
}

func TestHandleSearch_StateFilter(t *testing.T) {
	srv := newServerWith(t, filterFixture())
	resp, err := http.Get(srv.URL + "/search?state=Karnataka")
	require.NoError(t, err)
	defer resp.Body.Close()

	var body search.SearchResults
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Equal(t, 1, body.Total)
	assert.Equal(t, "SBIN0000002", body.Results[0].IFSC)
}

func TestHandleSearch_CityFilter_NoNaviMumbaiBleed(t *testing.T) {
	srv := newServerWith(t, filterFixture())
	resp, err := http.Get(srv.URL + "/search?city=Mumbai")
	require.NoError(t, err)
	defer resp.Body.Close()

	var body search.SearchResults
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Equal(t, 1, body.Total, "city=Mumbai must not match NAVI MUMBAI")
	assert.Equal(t, "HDFC0000001", body.Results[0].IFSC)
}

func TestHandleSearch_DistrictFilter_MultiWord(t *testing.T) {
	srv := newServerWith(t, filterFixture())
	resp, err := http.Get(srv.URL + "/search?district=Bangalore+Urban")
	require.NoError(t, err)
	defer resp.Body.Close()

	var body search.SearchResults
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	require.Equal(t, 1, body.Total)
	assert.Equal(t, "SBIN0000002", body.Results[0].IFSC)
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
	assert.Len(t, body, 1, "/healthz must return only {status: ok}")
}

func TestHandleStatus(t *testing.T) {
	srv := newTestServer(t)
	resp, err := http.Get(srv.URL + "/status")
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
	resp, err := http.Get(srv.URL + "/list")
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

	srv := httptest.NewServer(newRouter(s, search.Version{Tag: "test"}, "/ifsc", io.Discard))
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

func TestAccessLog_CombinedFormat(t *testing.T) {
	branches := []*search.Branch{
		{IFSC: "HDFC0000001", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "ANDHERI WEST", City: "MUMBAI", State: "MAHARASHTRA"},
	}
	s, err := search.NewMemorySearcher(branches)
	require.NoError(t, err)
	t.Cleanup(func() { _ = s.Close() })

	var buf bytes.Buffer
	srv := httptest.NewServer(newRouter(s, search.Version{Tag: "test"}, "", &buf))
	t.Cleanup(srv.Close)

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/ifsc/HDFC0000001", nil)
	require.NoError(t, err)
	req.Header.Set("User-Agent", "test-agent/1.0")
	req.Header.Set("Referer", "https://example.test/page")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	resp.Body.Close()

	line := buf.String()
	pattern := `^127\.0\.0\.1 - - \[\d{2}/[A-Z][a-z]{2}/\d{4}:\d{2}:\d{2}:\d{2} [-+]\d{4}\] ` +
		`"GET /ifsc/HDFC0000001 HTTP/1\.1" 200 \d+ ` +
		`"https://example\.test/page" "test-agent/1\.0"\n$`
	assert.Regexp(t, regexp.MustCompile(pattern), line, "log line: %q", line)
}

func TestAccessLog_404StatusAndDashes(t *testing.T) {
	branches := []*search.Branch{
		{IFSC: "HDFC0000001", BankCode: "HDFC", BankName: "HDFC Bank",
			Branch: "ANDHERI WEST", City: "MUMBAI", State: "MAHARASHTRA"},
	}
	s, err := search.NewMemorySearcher(branches)
	require.NoError(t, err)
	t.Cleanup(func() { _ = s.Close() })

	var buf bytes.Buffer
	srv := httptest.NewServer(newRouter(s, search.Version{Tag: "test"}, "", &buf))
	t.Cleanup(srv.Close)

	req, err := http.NewRequest(http.MethodGet, srv.URL+"/ifsc/ZZZZ0000000", nil)
	require.NoError(t, err)
	req.Header.Set("User-Agent", "")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	resp.Body.Close()

	line := buf.String()
	assert.Contains(t, line, `"GET /ifsc/ZZZZ0000000 HTTP/1.1" 404 `)
	assert.Contains(t, line, `"-" "-"`,
		"missing referer and user-agent both render as -")
}
