# syntax=docker/dockerfile:1.7
# Multi-stage build for bank-search.
# Stage 1 compiles the binaries.
# Stage 2 runs build-index to bake a Bleve index for a chosen razorpay/ifsc release.
# Stage 3 is a distroless runtime that ships only the server binary + index.

ARG GO_VERSION=1.23
ARG IFSC_TAG=""

FROM golang:${GO_VERSION}-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
        -o /out/bank-search . && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
        -o /out/build-index ./cmd/build-index

FROM build AS indexer
ARG IFSC_TAG
RUN /out/build-index -tag "${IFSC_TAG}" -out /index

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build               /out/bank-search /bank-search
COPY --from=indexer --chown=nonroot:nonroot /index /index
ENV BANK_SEARCH_INDEX_PATH=/index \
    BANK_SEARCH_PORT=8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/bank-search"]
