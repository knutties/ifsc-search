package search

import (
	"fmt"
	"sort"
	"strings"

	"github.com/blevesearch/bleve/v2"
)

// Bank is a distinct (bank_code, bank_name) pair surfaced by ListBanks.
type Bank struct {
	BankCode string `json:"bank_code"`
	BankName string `json:"bank_name"`
}

// ListBanks returns the distinct banks present in the index, sorted by
// bank_code. The result is computed once and cached for the lifetime of the
// searcher; the index is immutable at runtime so a single pass suffices.
func (b *bleveSearcher) ListBanks() ([]Bank, error) {
	b.banksOnce.Do(func() {
		b.banksCache, b.banksErr = b.computeBanks()
	})
	return b.banksCache, b.banksErr
}

func (b *bleveSearcher) computeBanks() ([]Bank, error) {
	q := bleve.NewMatchAllQuery()
	sr := bleve.NewSearchRequestOptions(q, 0, 0, false)
	// 10000 is well above the realistic count of distinct IFSC bank codes
	// (~200), so all terms come back in a single facet.
	sr.AddFacet("bank_code", bleve.NewFacetRequest("bank_code", 10000))

	res, err := b.idx.Search(sr)
	if err != nil {
		return nil, fmt.Errorf("bank facet: %w", err)
	}

	facet := res.Facets["bank_code"]
	if facet == nil || facet.Terms == nil {
		return []Bank{}, nil
	}
	terms := facet.Terms.Terms()
	banks := make([]Bank, 0, len(terms))
	for _, t := range terms {
		// The standard analyzer lowercases indexed terms; restore the
		// canonical upper-case bank code.
		code := strings.ToUpper(t.Term)
		name, err := b.bankNameFor(code)
		if err != nil {
			return nil, err
		}
		banks = append(banks, Bank{BankCode: code, BankName: name})
	}
	sort.Slice(banks, func(i, j int) bool { return banks[i].BankCode < banks[j].BankCode })
	return banks, nil
}

func (b *bleveSearcher) bankNameFor(code string) (string, error) {
	tq := bleve.NewTermQuery(strings.ToLower(code))
	tq.SetField("bank_code")
	sr := bleve.NewSearchRequestOptions(tq, 1, 0, false)
	sr.Fields = []string{"bank_name"}
	res, err := b.idx.Search(sr)
	if err != nil {
		return "", fmt.Errorf("bank name lookup for %q: %w", code, err)
	}
	if res.Total == 0 {
		return "", nil
	}
	name, _ := res.Hits[0].Fields["bank_name"].(string)
	return name, nil
}
