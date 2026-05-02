package search

import (
	"fmt"
	"sync"

	"github.com/blevesearch/bleve/v2"
)

// Searcher is the surface the HTTP layer depends on. Both the on-disk index
// and the in-memory test helper satisfy it.
type Searcher interface {
	Search(req SearchRequest) (*SearchResults, error)
	Lookup(code string) (*Branch, error)
	ListBanks() ([]Bank, error)
	DocCount() uint64
	Close() error
}

// bleveSearcher wraps a bleve.Index and adapts it to the Searcher interface.
// The actual Search method is implemented in query.go (Task 5) so this file
// only owns lifecycle.
type bleveSearcher struct {
	idx bleve.Index

	banksOnce  sync.Once
	banksCache []Bank
	banksErr   error
}

func newBleveSearcher(idx bleve.Index) *bleveSearcher {
	return &bleveSearcher{idx: idx}
}

func (b *bleveSearcher) DocCount() uint64 {
	n, err := b.idx.DocCount()
	if err != nil {
		return 0
	}
	return n
}

func (b *bleveSearcher) Close() error {
	return b.idx.Close()
}

// OpenIndex opens an existing on-disk Bleve index.
func OpenIndex(path string) (Searcher, error) {
	idx, err := bleve.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open bleve index at %q: %w", path, err)
	}
	return newBleveSearcher(idx), nil
}

// NewMemorySearcher builds an in-memory Bleve index seeded with the supplied
// branches. Used by tests so they do not need on-disk fixtures.
func NewMemorySearcher(branches []*Branch) (Searcher, error) {
	idx, err := bleve.NewMemOnly(NewIndexMapping())
	if err != nil {
		return nil, fmt.Errorf("create in-memory index: %w", err)
	}
	batch := idx.NewBatch()
	for _, b := range branches {
		if err := IndexBranch(batch, b); err != nil {
			return nil, fmt.Errorf("index %s: %w", b.IFSC, err)
		}
	}
	if err := idx.Batch(batch); err != nil {
		return nil, fmt.Errorf("commit batch: %w", err)
	}
	return newBleveSearcher(idx), nil
}

