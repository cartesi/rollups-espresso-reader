# (c) Cartesi and individual authors (see AUTHORS)
# SPDX-License-Identifier: Apache-2.0 (see LICENSE)

# syntax=docker.io/docker/dockerfile:1

ARG EMULATOR_VERSION=0.18.1

# Build directories.
ARG GO_BUILD_PATH=/build/cartesi/go

FROM cartesi/machine-emulator:${EMULATOR_VERSION} AS common-env

USER root

# Re-declare ARGs so they can be used in the RUN block
ARG GO_BUILD_PATH

# Install ca-certificates and curl (setup).
RUN <<EOF
    set -e
    apt-get update
    apt-get install -y --no-install-recommends ca-certificates curl wget build-essential pkg-config libssl-dev
    mkdir -p /opt/go ${GO_BUILD_PATH}/rollups-node
    chown -R cartesi:cartesi /opt/go ${GO_BUILD_PATH}
EOF

USER cartesi

# =============================================================================
# STAGE: go-installer
#
# This stage installs Go in the /opt directory.
# =============================================================================

FROM common-env AS go-installer
# Download and verify Go based on the target architecture
RUN <<EOF
    set -e
    ARCH=$(dpkg --print-architecture)
    wget -O /tmp/go.tar.gz "https://go.dev/dl/go1.22.7.linux-${ARCH}.tar.gz"
    sha256sum /tmp/go.tar.gz
    case "$ARCH" in
        amd64) echo "fc5d49b7a5035f1f1b265c17aa86e9819e6dc9af8260ad61430ee7fbe27881bb  /tmp/go.tar.gz" | sha256sum --check ;;
        arm64) echo "ed695684438facbd7e0f286c30b7bc2411cfc605516d8127dc25c62fe5b03885  /tmp/go.tar.gz" | sha256sum --check ;;
        *) echo "unsupported architecture: $ARCH"; exit 1 ;;
    esac
    tar -C /opt -xzf /tmp/go.tar.gz
    rm /tmp/go.tar.gz
EOF

# Set up Go environment variables
ENV PATH="/opt/go/bin:$PATH"

# =============================================================================
# STAGE: go-prepare
#
# This stage prepares the Go build environment. It downloads the external
# =============================================================================

FROM go-installer AS go-prepare

ARG GO_BUILD_PATH
WORKDIR ${GO_BUILD_PATH}

ENV GOCACHE=${GO_BUILD_PATH}/.cache
ENV GOENV=${GO_BUILD_PATH}/.config/go/env
ENV GOPATH=${GO_BUILD_PATH}/.go

# =============================================================================
# STAGE: go-builder
#
# This stage builds the node Go binaries. First it downloads the external
# dependencies and then it builds the binaries.
# =============================================================================

FROM go-prepare AS go-builder

ARG GO_BUILD_PATH

COPY --chown=cartesi:cartesi . ${GO_BUILD_PATH}/

# Build espresso reader.
RUN go build -o cartesi-rollups-espresso-reader

# =============================================================================
# STAGE: rollups-node
#
# This stage prepares the final Docker image that will be used in the production
# environment. It installs in /usr/bin all the binaries necessary to run the
# node.
#
# (This stage copies the binaries from previous stages.)
# =============================================================================

FROM debian:bookworm-20250113-slim AS rollups-node

ARG NODE_RUNTIME_DIR=/var/lib/cartesi-rollups-node

USER root

# Download system dependencies required at runtime.
ARG DEBIAN_FRONTEND=noninteractive
RUN <<EOF
    set -e
    apt-get update
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        libboost-all-dev \
        libslirp0 \
        lua5.4 \
        procps \
        tini
    rm -rf /var/lib/apt/lists/*
    mkdir -p ${NODE_RUNTIME_DIR}/snapshots ${NODE_RUNTIME_DIR}/data
EOF

# Copy Go binary.
ARG GO_BUILD_PATH

# Install cartesi machine
ENV CARTESI_MACHINE_VERSION=0.19.0-alpha4
RUN <<EOF
    set -e
    TARGETARCH=$(dpkg --print-architecture)
    curl -LO https://github.com/cartesi/machine-emulator/releases/download/v${CARTESI_MACHINE_VERSION}/machine-emulator_${TARGETARCH}.deb
    dpkg -i ./machine-emulator_${TARGETARCH}.deb
    rm ./machine-emulator_${TARGETARCH}.deb
EOF

RUN curl -L https://github.com/cartesi/image-kernel/releases/download/v0.20.0/linux-6.5.13-ctsi-1-v0.20.0.bin -o /usr/share/cartesi-machine/images/linux.bin

RUN curl -L https://github.com/cartesi/machine-emulator-tools/releases/download/v0.16.1/rootfs-tools-v0.16.1.ext2 -o /usr/share/cartesi-machine/images/rootfs.ext2

# Install rollups node
ARG ROLLUPS_NODE_VERSION
RUN <<EOF
    set -e
    TARGETARCH=$(dpkg --print-architecture)
    curl -LO https://github.com/cartesi/rollups-node/releases/download/v${ROLLUPS_NODE_VERSION}/cartesi-rollups-node-v${ROLLUPS_NODE_VERSION}_${TARGETARCH}.deb
    dpkg -i ./cartesi-rollups-node-v${ROLLUPS_NODE_VERSION}_${TARGETARCH}.deb
    rm ./cartesi-rollups-node-v${ROLLUPS_NODE_VERSION}_${TARGETARCH}.deb
EOF

WORKDIR ${NODE_RUNTIME_DIR}

# HEALTHCHECK --interval=1s --timeout=1s --retries=5 \
    # CMD curl -G -f -H 'Content-Type: application/json' http://127.0.0.1:10000/healthz

RUN mkdir applications
RUN cartesi-machine --ram-length=128Mi --store=applications/echo-dapp --final-hash -- ioctl-echo-loop --vouchers=1 --notices=1 --reports=1 --verbose=1

COPY --from=go-builder ${GO_BUILD_PATH}/cartesi-rollups-* /usr/bin

COPY --chown=cartesi:cartesi ./ci/entrypoint.sh /entrypoint.sh
COPY --chown=cartesi:cartesi ./ci/scripts/migrate-db.sh /migrate-db.sh
COPY --chown=cartesi:cartesi ./ci/scripts/espresso-dev-node-health.sh /espresso-dev-node-health.sh
COPY --chown=cartesi:cartesi ./internal/repository/postgres/schema/migrations /internal/repository/postgres/schema/migrations
ADD https://raw.githubusercontent.com/cartesi/rollups-node/refs/tags/v2.0.0-dev-20250527/internal/repository/postgres/schema/migrations/000001_create_initial_schema.up.sql /rollups-node/migrations/

RUN chmod +x /entrypoint.sh
CMD ["tini", "--", "/entrypoint.sh"]
