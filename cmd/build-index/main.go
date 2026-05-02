// Command build-index downloads IFSC.csv from a razorpay/ifsc GitHub release
// and produces a Bleve index plus a version.json sidecar.
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/knutties/ifsc-search/search"
)

const (
	githubAPI     = "https://api.github.com/repos/razorpay/ifsc/releases"
	batchSize     = 1000
	httpUserAgent = "ifsc-search-build-index/1.0"
)

func main() {
	tag := flag.String("tag", "", "release tag to use (default: latest)")
	indexDir := flag.String("out", "ifsc-api/index", "output index directory")
	csvPath := flag.String("csv", "", "path to a local IFSC.csv (skips download)")
	flag.Parse()

	if *csvPath == "" {
		downloaded, releaseTag, rbiDate, err := downloadCSV(*tag)
		if err != nil {
			log.Fatalf("download CSV: %v", err)
		}
		defer os.Remove(downloaded)
		*csvPath = downloaded
		*tag = releaseTag

		count, err := buildIndexFromCSV(*csvPath, *indexDir)
		if err != nil {
			log.Fatalf("build index: %v", err)
		}
		v := search.Version{
			Tag:           releaseTag,
			RBIUpdateDate: rbiDate,
			BuiltAt:       time.Now().UTC().Format(time.RFC3339),
		}
		if err := v.Save(*indexDir); err != nil {
			log.Fatalf("save version: %v", err)
		}
		log.Printf("built index for %s (%d docs)", releaseTag, count)
		return
	}

	count, err := buildIndexFromCSV(*csvPath, *indexDir)
	if err != nil {
		log.Fatalf("build index: %v", err)
	}
	v := search.Version{
		Tag:     *tag,
		BuiltAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := v.Save(*indexDir); err != nil {
		log.Fatalf("save version: %v", err)
	}
	log.Printf("built index from %s (%d docs)", *csvPath, count)
}

func buildIndexFromCSV(csvPath, indexDir string) (int, error) {
	if err := os.RemoveAll(indexDir); err != nil {
		return 0, fmt.Errorf("clean %s: %w", indexDir, err)
	}
	if err := os.MkdirAll(filepath.Dir(indexDir), 0755); err != nil {
		return 0, fmt.Errorf("mkdir parent: %w", err)
	}

	f, err := os.Open(csvPath)
	if err != nil {
		return 0, fmt.Errorf("open %s: %w", csvPath, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1 // tolerate trailing-comma quirks

	header, err := r.Read()
	if err != nil {
		return 0, fmt.Errorf("read header: %w", err)
	}
	cols, err := search.NewColumnIndex(header)
	if err != nil {
		return 0, err
	}

	idx, err := bleve.New(indexDir, search.NewIndexMapping())
	if err != nil {
		return 0, fmt.Errorf("create index: %w", err)
	}
	defer idx.Close()

	batch := idx.NewBatch()
	count := 0
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("read row: %w", err)
		}
		b, err := search.BranchFromCSVRow(cols, row)
		if err != nil {
			log.Printf("skipping row: %v", err)
			continue
		}
		if err := search.IndexBranch(batch, b); err != nil {
			return 0, fmt.Errorf("index %s: %w", b.IFSC, err)
		}
		count++
		if count%batchSize == 0 {
			if err := idx.Batch(batch); err != nil {
				return 0, fmt.Errorf("commit batch: %w", err)
			}
			batch = idx.NewBatch()
		}
	}
	if batch.Size() > 0 {
		if err := idx.Batch(batch); err != nil {
			return 0, fmt.Errorf("commit final batch: %w", err)
		}
	}
	return count, nil
}

type ghRelease struct {
	TagName string `json:"tag_name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name        string `json:"name"`
		DownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var rbiDateRe = regexp.MustCompile(`RBI Update Date.*?` + "`" + `([0-9-]+)` + "`")

func downloadCSV(tag string) (path, resolvedTag, rbiDate string, err error) {
	url := githubAPI + "/latest"
	if tag != "" {
		url = githubAPI + "/tags/" + tag
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", httpUserAgent)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", "", fmt.Errorf("github api: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", "", "", fmt.Errorf("github api status %d", resp.StatusCode)
	}
	var rel ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", "", "", fmt.Errorf("decode release: %w", err)
	}

	var dlURL string
	for _, a := range rel.Assets {
		if a.Name == "IFSC.csv" {
			dlURL = a.DownloadURL
			break
		}
	}
	if dlURL == "" {
		return "", "", "", fmt.Errorf("release %s has no IFSC.csv asset", rel.TagName)
	}

	tmp, err := os.CreateTemp("", "ifsc-*.csv")
	if err != nil {
		return "", "", "", err
	}
	defer tmp.Close()

	dlReq, _ := http.NewRequest("GET", dlURL, nil)
	dlReq.Header.Set("User-Agent", httpUserAgent)
	dlResp, err := http.DefaultClient.Do(dlReq)
	if err != nil {
		os.Remove(tmp.Name())
		return "", "", "", fmt.Errorf("download csv: %w", err)
	}
	defer dlResp.Body.Close()
	if dlResp.StatusCode != 200 {
		os.Remove(tmp.Name())
		return "", "", "", fmt.Errorf("download csv status %d", dlResp.StatusCode)
	}
	if _, err := io.Copy(tmp, dlResp.Body); err != nil {
		os.Remove(tmp.Name())
		return "", "", "", fmt.Errorf("copy csv: %w", err)
	}

	if m := rbiDateRe.FindStringSubmatch(rel.Body); len(m) == 2 {
		rbiDate = m[1]
	}
	return tmp.Name(), rel.TagName, rbiDate, nil
}
