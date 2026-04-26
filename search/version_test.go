package search

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadVersion_ReadsFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "version.json")
	require.NoError(t, os.WriteFile(path, []byte(
		`{"tag":"v2.0.57","rbi_update_date":"2026-03-15",`+
			`"indexed_docs":178670,"built_at":"2026-04-26T10:12:33Z"}`), 0644))

	v, err := LoadVersion(dir)
	require.NoError(t, err)
	assert.Equal(t, "v2.0.57", v.Tag)
	assert.Equal(t, "2026-03-15", v.RBIUpdateDate)
	assert.Equal(t, 178670, v.IndexedDocs)
	assert.Equal(t, "2026-04-26T10:12:33Z", v.BuiltAt)
}

func TestLoadVersion_MissingFile_ReturnsZeroValue(t *testing.T) {
	v, err := LoadVersion(t.TempDir())
	require.NoError(t, err)
	assert.Equal(t, "", v.Tag)
	assert.Equal(t, 0, v.IndexedDocs)
}
