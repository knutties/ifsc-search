# syntax=docker/dockerfile:1.7
# Multi-stage build for ifsc-search.
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
        -o /out/ifsc-search ./ifsc-api && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
        -o /out/build-index ./ifsc-api/cmd/build-index

FROM build AS indexer
ARG IFSC_TAG
RUN /out/build-index -tag "${IFSC_TAG}" -out /index

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build               /out/ifsc-search /ifsc-search
COPY --from=indexer --chown=nonroot:nonroot /index /index
ENV IFSC_SEARCH_INDEX_PATH=/index \
    IFSC_SEARCH_PORT=8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/ifsc-search"]
