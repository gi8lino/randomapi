# Build the manager binary
FROM golang:1.26 AS prep

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
ARG COMMIT=none
ARG LDFLAGS="-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT}"

ENV CGO_ENABLED=0

WORKDIR /workspace

# Copy the Go module manifests first so dependency downloads can be cached.
COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

# Copy the Go source.
COPY cmd/ cmd/
COPY internal/ internal

# Build the binary.
# TARGETARCH is intentionally allowed to be empty for regular Docker builds,
# so Go uses the builder container's default architecture.
RUN GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
  go build -ldflags="$LDFLAGS" -a -o randomapi ./cmd/main.go

RUN mkdir -p /outfs/work /outfs/tmp \
  # Change group ownership to GID 0 (root group),
  # because OpenShift assigns containers a random UID but keeps them compatible
  # with root-group-owned writable paths.
  && chgrp -R 0 /outfs/work /outfs/tmp \
  # Give group 0 read/write/execute permissions.
  # X only applies to directories or files that are already executable.
  && chmod -R g+rwX /outfs/work /outfs/tmp \
  # Set the setgid bit so new files/directories inherit group 0.
  && chmod g+s /outfs/work /outfs/tmp

# Use distroless as minimal base image to package the manager binary.
# Refer to https://github.com/GoogleContainerTools/distroless for more details.
FROM gcr.io/distroless/static:nonroot

COPY --from=prep /workspace/randomapi /randomapi
COPY --from=prep /outfs/work /work
COPY --from=prep /outfs/tmp /tmp

ENV HOME=/tmp
WORKDIR /work

# Run as a non-root user by default.
# Use GID 0 so the process can write to root-group-owned writable paths,
# which keeps the image compatible with OpenShift's arbitrary UID model.
USER 65532:0

ENTRYPOINT ["/randomapi"]
