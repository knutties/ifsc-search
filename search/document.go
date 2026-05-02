package search

import (
	"fmt"
	"strings"

	"github.com/blevesearch/bleve/v2"
)

// Branch is the indexed and returned representation of one IFSC record.
// Field tags match the JSON contract documented in the spec.
type Branch struct {
	IFSC     string `json:"ifsc"`
	BankCode string `json:"bank_code"`
	BankName string `json:"bank_name"`
	Branch   string `json:"branch"`
	Centre   string `json:"centre"`
	District string `json:"district"`
	State    string `json:"state"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Contact  string `json:"contact"`
	MICR     string `json:"micr"`
	SWIFT    string `json:"swift"`
	UPI      bool   `json:"upi"`
	NEFT     bool   `json:"neft"`
	RTGS     bool   `json:"rtgs"`
	IMPS     bool   `json:"imps"`
}

// requiredColumns lists CSV headers that must be present. ISO3166 is ignored
// because the spec does not surface country/state codes through the API.
var requiredColumns = []string{
	"BANK", "IFSC", "BRANCH", "CENTRE", "DISTRICT", "STATE", "ADDRESS",
	"CONTACT", "IMPS", "RTGS", "CITY", "NEFT", "MICR", "UPI", "SWIFT",
}

// ColumnIndex maps an IFSC.csv header name to its column index in the row.
type ColumnIndex map[string]int

func NewColumnIndex(header []string) (ColumnIndex, error) {
	idx := make(ColumnIndex, len(header))
	for i, h := range header {
		idx[strings.TrimSpace(h)] = i
	}
	for _, col := range requiredColumns {
		if _, ok := idx[col]; !ok {
			return nil, fmt.Errorf("CSV header missing required column %q", col)
		}
	}
	return idx, nil
}

func (c ColumnIndex) get(row []string, name string) string {
	i, ok := c[name]
	if !ok || i >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[i])
}

func BranchFromCSVRow(cols ColumnIndex, row []string) (*Branch, error) {
	ifsc := cols.get(row, "IFSC")
	if len(ifsc) < 4 {
		return nil, fmt.Errorf("IFSC %q too short", ifsc)
	}
	return &Branch{
		IFSC:     ifsc,
		BankCode: ifsc[0:4],
		BankName: cols.get(row, "BANK"),
		Branch:   cols.get(row, "BRANCH"),
		Centre:   cols.get(row, "CENTRE"),
		District: cols.get(row, "DISTRICT"),
		State:    cols.get(row, "STATE"),
		Address:  cols.get(row, "ADDRESS"),
		City:     cols.get(row, "CITY"),
		Contact:  cols.get(row, "CONTACT"),
		MICR:     cols.get(row, "MICR"),
		SWIFT:    cols.get(row, "SWIFT"),
		UPI:      parseBool(cols.get(row, "UPI")),
		NEFT:     parseBool(cols.get(row, "NEFT")),
		RTGS:     parseBool(cols.get(row, "RTGS")),
		IMPS:     parseBool(cols.get(row, "IMPS")),
	}, nil
}

func parseBool(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "true", "1", "yes", "y":
		return true
	default:
		return false
	}
}

// newIndexDoc returns a flat map containing every Branch field plus the
// lowercased keyword companion fields (state_key, district_key, city_key,
// ifsc_key) used by strict-equality filters. The map shape is used instead
// of an embedded struct because Bleve's reflection-based document walker
// does not promote anonymous embedded struct fields.
func newIndexDoc(b *Branch) map[string]interface{} {
	return map[string]interface{}{
		"ifsc":         b.IFSC,
		"bank_code":    b.BankCode,
		"bank_name":    b.BankName,
		"branch":       b.Branch,
		"centre":       b.Centre,
		"district":     b.District,
		"state":        b.State,
		"address":      b.Address,
		"city":         b.City,
		"contact":      b.Contact,
		"micr":         b.MICR,
		"swift":        b.SWIFT,
		"upi":          b.UPI,
		"neft":         b.NEFT,
		"rtgs":         b.RTGS,
		"imps":         b.IMPS,
		"state_key":    strings.ToLower(b.State),
		"district_key": strings.ToLower(b.District),
		"city_key":     strings.ToLower(b.City),
		"ifsc_key":     strings.ToLower(b.IFSC),
	}
}

// IndexBranch indexes b into target (a bleve.Index or bleve.Batch) using the
// IFSC code as the document id. Callers should prefer this helper over
// passing a bare *Branch so that strict-equality filter fields are
// populated.
func IndexBranch(target indexTarget, b *Branch) error {
	return target.Index(b.IFSC, newIndexDoc(b))
}

// indexTarget is satisfied by both bleve.Index and bleve.Batch — they share
// the Index(id, data) shape.
type indexTarget interface {
	Index(id string, data interface{}) error
}

// compile-time guards: bleve.Index and bleve.Batch must satisfy indexTarget.
var (
	_ indexTarget = (bleve.Index)(nil)
	_ indexTarget = (*bleve.Batch)(nil)
)
