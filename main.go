package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
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
		Handler:      newRouter(searcher, version, prefix, os.Stdout),
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

func newRouter(s search.Searcher, v search.Version, prefix string, accessLog io.Writer) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(prefix+"/search", handleSearch(s))
	mux.HandleFunc(prefix+"/healthz", handleHealthz())
	mux.HandleFunc(prefix+"/status", handleStatus(s, v))
	mux.HandleFunc("GET "+prefix+"/list", handleListBanks(s))
	mux.HandleFunc("GET "+prefix+"/ifsc/{code}", handleLookup(s))
	if accessLog == nil {
		return mux
	}
	return withAccessLog(mux, log.New(accessLog, "", 0))
}

// withAccessLog wraps next so each served request emits a single line in
// Apache Combined Log Format to logger:
//
//	%h - - [%t] "%r" %>s %b "%{Referer}i" "%{User-Agent}i"
func withAccessLog(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil || host == "" {
			host = r.RemoteAddr
		}
		if host == "" {
			host = "-"
		}
		ref := r.Header.Get("Referer")
		if ref == "" {
			ref = "-"
		}
		ua := r.Header.Get("User-Agent")
		if ua == "" {
			ua = "-"
		}
		logger.Printf(`%s - - [%s] "%s %s %s" %d %d "%s" "%s"`,
			host,
			start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method, r.URL.RequestURI(), r.Proto,
			rec.status, rec.bytes, ref, ua,
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status      int
	bytes       int
	wroteHeader bool
}

func (r *statusRecorder) WriteHeader(code int) {
	if r.wroteHeader {
		return
	}
	r.status = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	n, err := r.ResponseWriter.Write(b)
	r.bytes += n
	return n, err
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

func handleHealthz() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func handleStatus(s search.Searcher, v search.Version) http.HandlerFunc {
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
