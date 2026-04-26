package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const VersionFile = "version.json"

type Version struct {
	Tag           string `json:"tag"`
	RBIUpdateDate string `json:"rbi_update_date"`
	IndexedDocs   int    `json:"indexed_docs"`
	BuiltAt       string `json:"built_at"`
}

// LoadVersion reads {indexDir}/version.json. A missing file returns the zero
// value with no error so /healthz can still report something useful when the
// index was opened but the metadata file is absent.
func LoadVersion(indexDir string) (Version, error) {
	path := filepath.Join(indexDir, VersionFile)
	bytes, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return Version{}, nil
	}
	if err != nil {
		return Version{}, fmt.Errorf("read %s: %w", path, err)
	}
	var v Version
	if err := json.Unmarshal(bytes, &v); err != nil {
		return Version{}, fmt.Errorf("parse %s: %w", path, err)
	}
	return v, nil
}

// Save writes the version metadata next to the index. Used by build-index.
func (v Version) Save(indexDir string) error {
	path := filepath.Join(indexDir, VersionFile)
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal version: %w", err)
	}
	if err := os.WriteFile(path, bytes, 0644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
