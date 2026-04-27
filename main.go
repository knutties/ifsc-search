package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/knutties/ifsc-search/search"
)

func main() {
	port := envOr("IFSC_SEARCH_PORT", "8080")
	indexPath := envOr("IFSC_SEARCH_INDEX_PATH", "./ifsc-api/index")
	prefix := normalizePrefix(os.Getenv("PATH_PREFIX"))

	searcher, err := search.OpenIndex(indexPath)
	if err != nil {
		log.Fatalf("open index: %v", err)
	}
	defer searcher.Close()

	version, err := search.LoadVersion(indexPath)
	if err != nil {
		log.Printf("warning: could not load version metadata: %v", err)
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      newRouter(searcher, version, prefix),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		log.Printf("ifsc-search listening on :%s (index=%s, docs=%d, prefix=%q)",
			port, indexPath, searcher.DocCount(), prefix)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Println("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func newRouter(s search.Searcher, v search.Version, prefix string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(prefix+"/search", handleSearch(s))
	mux.HandleFunc(prefix+"/healthz", handleHealthz(s, v))
	mux.HandleFunc("GET "+prefix+"/banks", handleListBanks(s))
	mux.HandleFunc("GET "+prefix+"/ifsc/{code}", handleLookup(s))
	return mux
}

// normalizePrefix returns "" for empty input, otherwise ensures a single
// leading "/" and no trailing "/". e.g. "ifsc/" -> "/ifsc", "/" -> "".
func normalizePrefix(p string) string {
	p = strings.TrimSpace(p)
	p = strings.Trim(p, "/")
	if p == "" {
		return ""
	}
	return "/" + p
}

func handleSearch(s search.Searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		req, err := parseRequest(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		res, err := s.Search(req)
		if err != nil {
			if errors.Is(err, search.ErrMissingQuery) ||
				errors.Is(err, search.ErrBadPagination) {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			log.Printf("search error: %v", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		writeJSON(w, http.StatusOK, res)
	}
}

func parseRequest(r *http.Request) (search.SearchRequest, error) {
	q := r.URL.Query()
	req := search.SearchRequest{
		Bank: q.Get("bank"),
		Q:    q.Get("q"),
	}
	if v := q.Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return req, errors.New("limit must be an integer")
		}
		req.Limit = n
	}
	if v := q.Get("offset"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return req, errors.New("offset must be an integer")
		}
		req.Offset = n
	}
	return req, nil
}

func handleLookup(s search.Searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		br, err := s.Lookup(code)
		if err != nil {
			if errors.Is(err, search.ErrNotFound) {
				writeError(w, http.StatusNotFound, "ifsc code not found")
				return
			}
			log.Printf("lookup error: %v", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		writeJSON(w, http.StatusOK, br)
	}
}

func handleListBanks(s search.Searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		banks, err := s.ListBanks()
		if err != nil {
			log.Printf("list banks error: %v", err)
			writeError(w, http.StatusInternalServerError, "internal error")
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"total": len(banks),
			"banks": banks,
		})
	}
}

func handleHealthz(s search.Searcher, v search.Version) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"status":          "ok",
			"indexed_docs":    s.DocCount(),
			"release_tag":     v.Tag,
			"rbi_update_date": v.RBIUpdateDate,
			"built_at":        v.BuiltAt,
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
