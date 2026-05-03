.PHONY: build-index build-index-from run test clean \
        smithy-build smithy-publish smithy-updates smithy-clean

build-index:
	go run ./cmd/build-index -out ./index

build-index-from:
	@if [ -z "$(CSV)" ]; then echo "usage: make build-index-from CSV=path/to/IFSC.csv"; exit 1; fi
	go run ./cmd/build-index -csv $(CSV) -out ./index

run:
	@if [ ! -d ./index ]; then \
	  echo "no index at ./index — run 'make build-index' (or 'make build-index-from CSV=...') first"; \
	  exit 1; \
	fi
	go run .

test:
	go test -tags=unit ./...

clean:
	rm -rf ./index

# Builds Smithy projections into smithy/build/.
smithy-build:
	cd smithy && smithy build

# Copies the generated OpenAPI spec, Smithy AST, and TypeScript client
# sources into api/ and clients/ so consumers don't need a smithy-cli.
# The rsync excludes preserve any local node_modules/dist artifacts that
# devs produce by running `npm install` / `npm run build` inside
# clients/typescript/.
smithy-publish: smithy-build
	cp smithy/build/smithy/source/openapi/BankSearch.openapi.json api/openapi.json
	cp smithy/build/smithy/source/model/model.json api/model.json
	mkdir -p clients/typescript
	rsync -a --delete \
	  --exclude='node_modules' \
	  --exclude='dist-cjs' \
	  --exclude='dist-es' \
	  --exclude='dist-types' \
	  --exclude='*.tsbuildinfo' \
	  --exclude='package-lock.json' \
	  --exclude='README.md' \
	  smithy/build/smithy/source/typescript-client-codegen/ clients/typescript/
	cp LICENSE clients/typescript/LICENSE
	jq '.license = "MIT"' clients/typescript/package.json > clients/typescript/package.json.tmp \
	  && mv clients/typescript/package.json.tmp clients/typescript/package.json
	node smithy/scripts/post-process.mjs

smithy-clean:
	rm -rf smithy/build

# Used by CI to assert the committed api/ and clients/ artifacts match
# the IDL.
smithy-updates: smithy-clean smithy-publish
