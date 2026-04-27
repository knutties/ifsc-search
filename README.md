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

```bash
make run
# IFSC_SEARCH_PORT and IFSC_SEARCH_INDEX_PATH override defaults
# PATH_PREFIX mounts all routes under a sub-path, e.g. PATH_PREFIX=/ifsc
# exposes /ifsc/search, /ifsc/healthz, /ifsc/banks, and /ifsc/ifsc/{code}
```

## API

### `GET /search`

Query params:

| Name     | Required        | Notes                                          |
| -------- | --------------- | ---------------------------------------------- |
| `bank`   | one of bank/q   | 4-char IFSC bank code or fuzzy bank name       |
| `q`      | one of bank/q   | free-text over branch, city, address           |
| `limit`  | no              | default 20, max 100                            |
| `offset` | no              | default 0                                      |

Example:

```bash
curl 'http://localhost:8080/search?bank=HDFC&q=andheri+west&limit=5'
```

### `GET /banks`

Returns the distinct list of banks present in the index, sorted by
`bank_code`. The list is computed once on first request and cached.

```bash
curl 'http://localhost:8080/banks'
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

### `GET /healthz`

Returns the index version metadata and document count.

## Tests

```bash
make test
```

## Container image

A multi-stage `Dockerfile` at the repo root builds a distroless image that
ships `ifsc-search` plus a pre-baked Bleve index. The
`.github/workflows/image.yml` workflow_dispatch publishes images to
`ghcr.io/knutties/ifsc-search`.
