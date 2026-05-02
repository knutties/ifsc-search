package search

import (
	"errors"
	"fmt"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/search/query"
)

const (
	defaultLimit = 20
	maxLimit     = 100
	// bankNameMatchFloor is the minimum Bleve score a fuzzy bank-name match
	// must achieve to be accepted. Tuned empirically; see the spec.
	bankNameMatchFloor = 0.1
)

var (
	ErrMissingQuery  = errors.New("at least one of bank, q, ifsc, state, district, city is required")
	ErrBadPagination = errors.New("invalid pagination")
)

type SearchRequest struct {
	Bank       string
	Q          string
	IFSCPrefix string
	State      string
	District   string
	City       string
	Limit      int
	Offset     int
}

func (r *SearchRequest) hasSignal() bool {
	for _, v := range []string{r.Bank, r.Q, r.IFSCPrefix, r.State, r.District, r.City} {
		if strings.TrimSpace(v) != "" {
			return true
		}
	}
	return false
}

func (r *SearchRequest) Validate() error {
	if !r.hasSignal() {
		return ErrMissingQuery
	}
	if r.Offset < 0 {
		return fmt.Errorf("%w: offset must be >= 0", ErrBadPagination)
	}
	if r.Limit < 0 {
		return fmt.Errorf("%w: limit must be >= 0", ErrBadPagination)
	}
	return nil
}

func (r *SearchRequest) normalize() {
	if r.Limit <= 0 {
		r.Limit = defaultLimit
	}
	if r.Limit > maxLimit {
		r.Limit = maxLimit
	}
}

type ResultItem struct {
	*Branch
	Score float64 `json:"score"`
}

type SearchResults struct {
	Total   int          `json:"total"`
	Limit   int          `json:"limit"`
	Offset  int          `json:"offset"`
	Results []ResultItem `json:"results"`
}

// Search implements Searcher.Search.
func (b *bleveSearcher) Search(req SearchRequest) (*SearchResults, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	req.normalize()

	bankCode, ok, err := b.resolveBank(req.Bank)
	if err != nil {
		return nil, err
	}
	if req.Bank != "" && !ok {
		// User named a bank we cannot resolve — empty results, not an error.
		return &SearchResults{
			Total:   0,
			Limit:   req.Limit,
			Offset:  req.Offset,
			Results: []ResultItem{},
		}, nil
	}

	q := buildQuery(bankCode, req)

	sr := bleve.NewSearchRequestOptions(q, req.Limit, req.Offset, false)
	sr.Fields = []string{"*"}
	// Stable alpha sort whenever there is no free-text query — the new
	// strict filters are equivalent to "narrow then list" and benefit from
	// deterministic ordering.
	if strings.TrimSpace(req.Q) == "" {
		sr.SortBy([]string{"branch"})
	}

	res, err := b.idx.Search(sr)
	if err != nil {
		return nil, fmt.Errorf("bleve search: %w", err)
	}

	out := &SearchResults{
		Total:   int(res.Total),
		Limit:   req.Limit,
		Offset:  req.Offset,
		Results: make([]ResultItem, 0, len(res.Hits)),
	}
	for _, hit := range res.Hits {
		br := branchFromFields(hit.Fields)
		br.IFSC = hit.ID
		out.Results = append(out.Results, ResultItem{Branch: br, Score: hit.Score})
	}
	return out, nil
}

// resolveBank returns the canonical 4-char bank code for the user's input.
// ok=false means no acceptable match was found and the caller should return
// empty results.
func (b *bleveSearcher) resolveBank(input string) (string, bool, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", true, nil
	}
	if len(input) == 4 {
		// Try as exact bank code first. The standard analyzer lowercases
		// indexed terms, so the term query input must be lowercase even
		// though we return the canonical upper-case form.
		exact := bleve.NewTermQuery(strings.ToLower(input))
		exact.SetField("bank_code")
		req := bleve.NewSearchRequestOptions(exact, 1, 0, false)
		res, err := b.idx.Search(req)
		if err != nil {
			return "", false, fmt.Errorf("bank exact lookup: %w", err)
		}
		if res.Total > 0 {
			return strings.ToUpper(input), true, nil
		}
	}

	// Fall back to fuzzy name resolution.
	fz := bleve.NewMatchQuery(input)
	fz.SetField("bank_name")
	fz.SetFuzziness(1)
	req := bleve.NewSearchRequestOptions(fz, 1, 0, false)
	req.Fields = []string{"bank_code"}
	res, err := b.idx.Search(req)
	if err != nil {
		return "", false, fmt.Errorf("bank fuzzy lookup: %w", err)
	}
	if res.Total == 0 || res.Hits[0].Score < bankNameMatchFloor {
		return "", false, nil
	}
	code, _ := res.Hits[0].Fields["bank_code"].(string)
	if code == "" {
		return "", false, nil
	}
	return code, true, nil
}

func buildQuery(bankCode string, req SearchRequest) query.Query {
	conj := bleve.NewConjunctionQuery()

	if bankCode != "" {
		// bank_code is indexed via the standard analyzer (lowercased),
		// so the term query needs the lowercase form.
		bq := bleve.NewTermQuery(strings.ToLower(bankCode))
		bq.SetField("bank_code")
		conj.AddQuery(bq)
	}

	if pfx := strings.TrimSpace(req.IFSCPrefix); pfx != "" {
		pq := bleve.NewPrefixQuery(strings.ToLower(pfx))
		pq.SetField("ifsc_key")
		conj.AddQuery(pq)
	}

	for _, f := range []struct {
		field, value string
	}{
		{"state_key", req.State},
		{"district_key", req.District},
		{"city_key", req.City},
	} {
		v := strings.TrimSpace(f.value)
		if v == "" {
			continue
		}
		tq := bleve.NewTermQuery(strings.ToLower(v))
		tq.SetField(f.field)
		conj.AddQuery(tq)
	}

	if q := strings.TrimSpace(req.Q); q != "" {
		conj.AddQuery(textQuery(q))
	}

	// Internal invariant: Search() validates before calling buildQuery, so
	// at least one clause must be present. A bare conjunction with zero
	// clauses would otherwise match every document.
	if len(conj.Conjuncts) == 0 {
		panic("buildQuery called with no clauses — Validate() bypassed")
	}
	return conj
}

func textQuery(q string) query.Query {
	tokens := strings.Fields(strings.ToLower(q))
	disj := bleve.NewDisjunctionQuery()
	for _, tok := range tokens {
		for _, field := range []string{"branch", "city", "address"} {
			boost := fieldBoosts[field]

			fq := bleve.NewFuzzyQuery(tok)
			fq.SetField(field)
			fq.SetFuzziness(2)
			fq.SetBoost(boost)
			disj.AddQuery(fq)

			mq := bleve.NewMatchQuery(tok)
			mq.SetField(field)
			mq.SetBoost(boost * 2) // exact match outranks fuzzy
			disj.AddQuery(mq)
		}

		// Treat tokens that look like an IFSC (or its prefix) as a hit
		// against the doc's IFSC. Lets a user paste a code into the
		// generic q box and find the branch without knowing about the
		// dedicated `ifsc` param.
		ipq := bleve.NewPrefixQuery(tok)
		ipq.SetField("ifsc_key")
		ipq.SetBoost(4.0) // outrank fuzzy text matches
		disj.AddQuery(ipq)
	}
	return disj
}

func branchFromFields(f map[string]interface{}) *Branch {
	get := func(name string) string {
		if v, ok := f[name].(string); ok {
			return v
		}
		return ""
	}
	getBool := func(name string) bool {
		if v, ok := f[name].(bool); ok {
			return v
		}
		return false
	}
	return &Branch{
		BankCode: get("bank_code"),
		BankName: get("bank_name"),
		Branch:   get("branch"),
		Centre:   get("centre"),
		District: get("district"),
		State:    get("state"),
		Address:  get("address"),
		City:     get("city"),
		Contact:  get("contact"),
		MICR:     get("micr"),
		SWIFT:    get("swift"),
		UPI:      getBool("upi"),
		NEFT:     getBool("neft"),
		RTGS:     getBool("rtgs"),
		IMPS:     getBool("imps"),
	}
}
