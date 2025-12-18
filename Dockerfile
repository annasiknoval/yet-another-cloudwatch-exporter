# syntax=docker/dockerfile:1

############################
# 1️⃣ Builder stage
############################
FROM docker.io/golang:1.24.0 AS builder
WORKDIR /src

# Copy go.mod and go.sum first (cache layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source
COPY . .

# These args are automatically provided by BuildKit when using buildx
ARG TARGETOS
ARG TARGETARCH

# Ensure reproducible, static build
ENV CGO_ENABLED=0

# Build YACE binary dynamically for the target platform
RUN GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags="-s -w" -o /out/yace ./cmd/yace

############################
# 2️⃣ Final image
############################
FROM quay.io/prometheus/busybox-linux-${TARGETARCH:-amd64}:latest

LABEL maintainer="The Prometheus Authors <prometheus-developers@googlegroups.com>"

# Copy binary and example config
COPY --from=builder /out/yace /bin/yace
COPY examples/ec2.yml /etc/yace/config.yml

EXPOSE 5000
USER nobody

ENTRYPOINT ["/bin/yace"]
CMD ["--config.file=/etc/yace/config.yml"]
