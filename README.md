# ifsc-search

A self-contained HTTP service for locating Indian bank branches by bank plus a
fuzzy free-text query over branch name, address, and city. Built on a Bleve
index generated from `IFSC.csv` shipped with each
[`razorpay/ifsc` release](https://github.com/razorpay/ifsc/releases).

## Build the index

```bash
make build-index                      # downloads latest release CSV
make build-index-from CSV=./IFSC.csv  # uses a local CSV
```

The index lands in `./index/` (gitignored) along with `version.json`.

## Run the server

### From source

Builds the binary locally and serves the index in `./index/`:

```bash
make run
# IFSC_SEARCH_PORT and IFSC_SEARCH_INDEX_PATH override defaults
# PATH_PREFIX mounts all routes under a sub-path, e.g. PATH_PREFIX=/ifsc
# exposes /ifsc/search, /ifsc/healthz, /ifsc/status, /ifsc/list, and /ifsc/ifsc/{code}
```

### From the published GHCR image

Pulls the distroless image from `ghcr.io/knutties/ifsc-search`, which already
ships a pre-baked Bleve index — no `make build-index` step required:

```bash
docker run --rm -p 8080:8080 ghcr.io/knutties/ifsc-search:latest
curl 'http://localhost:8080/healthz'
```

Mount under a sub-path with `PATH_PREFIX`:

```bash
docker run --rm -p 8080:8080 -e PATH_PREFIX=/ifsc \
    ghcr.io/knutties/ifsc-search:latest
```

Build the same image locally instead of pulling:

```bash
docker build -t ifsc-search .
docker run --rm -p 8080:8080 ifsc-search
```

## API

### `GET /search`

Query params:

| Name       | Required | Notes                                                   |
| ---------- | -------- | ------------------------------------------------------- |
| `bank`     | one of * | 4-char IFSC bank code or fuzzy bank name                |
| `q`        | one of * | free-text over branch, city, address, IFSC prefix       |
| `ifsc`     | one of * | case-insensitive IFSC prefix, e.g. `HDFC0CAG`           |
| `state`    | one of * | strict, case-insensitive (`Maharashtra`, `West Bengal`) |
| `district` | one of * | strict, case-insensitive                                |
| `city`     | one of * | strict, case-insensitive (no substring false-positives) |
| `limit`    | no       | default 20, max 100                                     |
| `offset`   | no       | default 0                                               |

\* at least one of `bank`, `q`, `ifsc`, `state`, `district`, `city` is
required. When more than one is supplied they AND-combine: e.g.
`bank=HDFC&state=Maharashtra&q=andheri` narrows to HDFC branches in
Maharashtra matching "andheri".

Strict filters (`state`, `district`, `city`) are exact-match against the
indexed value — `city=Mumbai` will not bleed into "Navi Mumbai".

Example:

```bash
curl 'http://localhost:8080/search?bank=HDFC&q=andheri+west&limit=5'
curl 'http://localhost:8080/search?ifsc=HDFC0CAG'
curl 'http://localhost:8080/search?bank=HDFC&city=Mumbai'
```

### `GET /list`

Returns the distinct list of banks present in the index, sorted by
`bank_code`. The list is computed once on first request and cached.

```bash
curl 'http://localhost:8080/list'
```

```json
{
  "total": 168,
  "banks": [
    {"bank_code": "ABHY", "bank_name": "Abhyudaya Co-operative Bank Limited"},
    {"bank_code": "HDFC", "bank_name": "HDFC Bank"}
  ]
}
```

### `GET /ifsc/{code}`

Returns the full branch JSON for a given IFSC code, or 404 if unknown.
Mirrors the URL shape of the hosted `https://ifsc.razorpay.com/{code}` API.

```bash
curl 'http://localhost:8080/ifsc/HDFC0CAGSBK'
```

### Access logs

Every served request emits one line on stdout in Apache Combined Log
Format, e.g.:

```
::1 - - [30/Apr/2026:19:55:08 +0530] "GET /ifsc/HDFC0000001 HTTP/1.1" 200 280 "https://example.test/" "curl/8.7.1"
```

### `GET /healthz`

Lightweight liveness probe. Returns `{"status": "ok"}` and nothing else —
safe to wire into load-balancer health checks without leaking build metadata.

```bash
curl 'http://localhost:8080/healthz'
```

### `GET /status`

Returns the index version metadata and document count.

```bash
curl 'http://localhost:8080/status'
```

## Tests

```bash
make test
```

## Container image

A multi-stage `Dockerfile` at the repo root builds a distroless image that
ships `ifsc-search` plus a pre-baked Bleve index. The
`.github/workflows/image.yml` workflow_dispatch publishes images to
`ghcr.io/knutties/ifsc-search`. See [Run the server](#run-the-server) for
usage.
