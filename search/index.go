// Package search owns the Bleve index used by the ifsc-search HTTP service.
package search

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

const DocType = "branch"

// fieldBoosts records per-field boosts so tests (and future tuning) can read
// them without poking at Bleve internals.
var fieldBoosts = map[string]float64{
	"branch":  3.0,
	"city":    2.0,
	"address": 1.0,
}

// FieldBoost returns the boost recorded for a field mapping. Returns 1.0 if
// the field is not in the boost table.
func FieldBoost(f *mapping.FieldMapping) float64 {
	if f == nil {
		return 1.0
	}
	if b, ok := fieldBoosts[f.Name]; ok {
		return b
	}
	return 1.0
}

func NewIndexMapping() *mapping.IndexMappingImpl {
	im := bleve.NewIndexMapping()
	im.DefaultAnalyzer = "standard"
	// Untyped documents fall back to the "branch" mapping declared below.
	// Without this, the registered field mappings (notably the keyword
	// analyzer pinned on *_key fields) are silently ignored and Bleve
	// auto-detects everything via the standard analyzer.
	im.DefaultType = DocType

	branch := bleve.NewDocumentMapping()

	for _, name := range []string{"bank_code", "bank_name", "branch", "city",
		"address", "centre", "district", "state"} {
		f := bleve.NewTextFieldMapping()
		f.Name = name
		f.Store = true
		f.Index = true
		// Boosts are applied at query time in query.go, not in the mapping.
		branch.AddFieldMappingsAt(name, f)
	}

	for _, name := range []string{"contact", "micr", "swift"} {
		f := bleve.NewKeywordFieldMapping()
		f.Name = name
		f.Store = true
		f.Index = false
		branch.AddFieldMappingsAt(name, f)
	}

	// Lowercased keyword companion fields used by strict-equality filters
	// (state/district/city) and the IFSC prefix filter. Populated by the
	// indexDoc wrapper so the public Branch JSON contract stays unchanged.
	// Analyzer is pinned to "keyword" so multi-word values like
	// "west bengal" stay as a single indexed term rather than being
	// tokenized by the index-level "standard" default.
	for _, name := range []string{"state_key", "district_key", "city_key", "ifsc_key"} {
		f := bleve.NewKeywordFieldMapping()
		f.Name = name
		f.Analyzer = "keyword"
		f.Store = false
		f.Index = true
		branch.AddFieldMappingsAt(name, f)
	}

	for _, name := range []string{"upi", "neft", "rtgs", "imps"} {
		f := bleve.NewBooleanFieldMapping()
		f.Name = name
		f.Store = true
		f.Index = false
		branch.AddFieldMappingsAt(name, f)
	}

	im.AddDocumentMapping(DocType, branch)
	return im
}
