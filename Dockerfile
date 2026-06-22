# syntax=docker/dockerfile:1.25
FROM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
ARG COMMIT=none
ARG LDFLAGS="-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT}"
ENV CGO_ENABLED=0

WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/randomapi/main.go cmd/randomapi/main.go
COPY internal/ internal

RUN GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -ldflags="${LDFLAGS}" -o /out/randomapi ./cmd/randomapi/main.go

RUN mkdir -p /outfs/work /outfs/tmp \
  # Change group ownership of /work and /tmp to GID 0 (root group),
  # because OpenShift assigns containers a random UID but always includes them in group 0.
  && chgrp -R 0 /outfs/work /outfs/tmp \
  # Give group 0 read/write/execute (X only applies to dirs or already-executable files).
  # This makes the dirs writable by arbitrary UIDs in group 0.
  && chmod -R g+rwX /outfs/work /outfs/tmp \
  # Set the setgid bit on the dirs so that any new files/dirs created inside
  # will inherit group 0 instead of the creator's primary group.
  && chmod g+s /outfs/work /outfs/tmp

FROM gcr.io/distroless/static:nonroot
COPY --from=builder /out/randomapi /randomapi
COPY --from=builder /outfs/work /work
COPY --from=builder /outfs/tmp  /tmp
ENV HOME=/tmp
WORKDIR /work
USER 65532:65532
ENTRYPOINT ["/randomapi"]

