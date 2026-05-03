# bank-search

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
# BANK_SEARCH_PORT and BANK_SEARCH_INDEX_PATH override defaults
# PATH_PREFIX mounts all routes under a sub-path, e.g. PATH_PREFIX=/ifsc
# exposes /ifsc/search, /ifsc/healthz, /ifsc/status, /ifsc/list, and /ifsc/ifsc/{ifsc}
```

### From the published GHCR image

Pulls the distroless image from `ghcr.io/knutties/bank-search`, which already
ships a pre-baked Bleve index — no `make build-index` step required:

```bash
docker run --rm -p 8080:8080 ghcr.io/knutties/bank-search:latest
curl 'http://localhost:8080/healthz'
```

Mount under a sub-path with `PATH_PREFIX`:

```bash
docker run --rm -p 8080:8080 -e PATH_PREFIX=/ifsc \
    ghcr.io/knutties/bank-search:latest
```

Build the same image locally instead of pulling:

```bash
docker build -t bank-search .
docker run --rm -p 8080:8080 bank-search
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

### `GET /ifsc/{ifsc}`

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

`go test` covers the Go side. End-to-end coverage lives in
`e2e/typescript/smoke.test.ts` and drives a running server through the
generated TypeScript SDK — exercising serialization, error
discrimination, and every HTTP route in one pass:

```bash
# in one shell
make build-index-from CSV=./cmd/build-index/testdata/sample.csv
make run

# in another
cd clients/typescript && npm install && npm run build
cd ../../e2e/typescript && npm install && npm test
```

CI's `e2e` job runs this same flow on every PR.

## API contract (Smithy)

The HTTP contract is modelled in [Smithy 2.0](https://smithy.io) under
`smithy/models/` and projected into a committed OpenAPI 3.1 spec at
`api/openapi.json`, a Smithy AST at `api/model.json`, and a TypeScript
client at `clients/typescript/`.

```bash
make smithy-build       # produces smithy/build/ (gitignored)
make smithy-publish     # also copies into api/ and clients/typescript/
make smithy-updates     # clean + publish; CI runs this and gates on no diff
```

The Go server stays hand-written; Smithy is the source of truth for the
wire format and drives every published client.

## TypeScript client

`clients/typescript/` ships an AWS-SDK-v3-style client generated by
`smithy-aws-typescript-codegen` against the same IDL. The package is
`@knutties/bank-search-client`. Install it from a local checkout:

```bash
cd clients/typescript && npm install && npm run build
```

```ts
import { BankSearchClient, GetBranchCommand, SearchCommand } from
    "@knutties/bank-search-client";

const client = new BankSearchClient({ endpoint: "http://localhost:8080" });

const branch = await client.send(new GetBranchCommand({ ifsc: "HDFC0000001" }));
const hits   = await client.send(new SearchCommand({ bank: "HDFC", q: "andheri" }));
```

The generated source is committed; `dist-*`, `node_modules`, and the lock
file are gitignored. Re-run `make smithy-updates` to refresh.

## Development environment

A `flake.nix` ships a reproducible dev shell with Go, `gopls`,
`golangci-lint`, `gh`, `jq`, `smithy-cli`, and `nodejs_22` pinned via
`nixpkgs`:

```bash
nix develop          # one-off shell
direnv allow         # auto-activate via the included .envrc
```

The flake pins Go 1.25 (currently supported upstream); the project's
`go.mod` still declares `go 1.23` for compatibility with the CI image
and Dockerfile. Newer toolchains build older modules without changes.

## Container image

A multi-stage `Dockerfile` at the repo root builds a distroless image that
ships `bank-search` plus a pre-baked Bleve index. The
`.github/workflows/image.yml` workflow_dispatch publishes images to
`ghcr.io/knutties/bank-search`. See [Run the server](#run-the-server) for
usage.
