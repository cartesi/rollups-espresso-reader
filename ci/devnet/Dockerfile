ARG FOUNDRY_VERSION=1.0.0
ARG CONTRACTS_VERSION=2.0.0-rc.18
ARG DEVNET_BUILD_PATH=/opt/cartesi/rollups-contracts

FROM debian:bookworm-20250407 AS rollups-node-devnet
ARG FOUNDRY_VERSION
ARG CONTRACTS_VERSION
ARG DEVNET_BUILD_PATH

USER root

# Install ca-certificates, curl, and git (setup).
ENV DEBIAN_FRONTEND=noninteractive
RUN <<EOF
    set -e
    apt-get update
    apt-get install -y --no-install-recommends ca-certificates curl git jq
EOF

# Install Foundry from downloaded pre-compiled binaries.
RUN <<EOF
    set -e
    URL=https://github.com/foundry-rs/foundry/releases/download
    ARCH=$(dpkg --print-architecture)
    ARTIFACT=foundry_v${FOUNDRY_VERSION}_linux_${ARCH}.tar.gz
    curl -sSL ${URL}/v${FOUNDRY_VERSION}/${ARTIFACT} -o /tmp/foundry.tar.gz
    case $ARCH in
    amd64) echo "8c078eaced6bfce76af902def65e36db677842ffbf53f76730647082a5b1e45a  /tmp/foundry.tar.gz" | shasum -a 256 --check ;;
    arm64) echo "79967815024de284265b9e82b75498b46357e597b93b481b9f4abe88a53eade6  /tmp/foundry.tar.gz" | shasum -a 256 --check ;;
    esac
    tar -zxf /tmp/foundry.tar.gz -C /usr/local/bin
EOF

# Install anvil state file and devnet deployment info.
RUN <<EOF
    set -e
    mkdir -p ${DEVNET_BUILD_PATH}
    URL=https://github.com/cartesi/rollups-contracts/releases/download
    VERSION=v${CONTRACTS_VERSION}
    ARTIFACT=rollups-contracts-${CONTRACTS_VERSION}-anvil-v${FOUNDRY_VERSION}.tar.gz
    curl -sSL ${URL}/${VERSION}/${ARTIFACT} -o /tmp/contracts.tar.gz
    echo "8962909ea5be1e190aa099a08a949710d3fe2810ccdadcd3b860ef4338f87ece  /tmp/contracts.tar.gz" | sha256sum --check
    tar -zxf /tmp/contracts.tar.gz -C ${DEVNET_BUILD_PATH}/
    mkdir -p /usr/share/devnet
    cat ${DEVNET_BUILD_PATH}/deployments/*.json | jq -s 'map({ (.contractName): .address }) | add' > /usr/share/devnet/deployment.json
    mv ${DEVNET_BUILD_PATH}/state.json /usr/share/devnet/anvil_state.json
EOF

HEALTHCHECK --interval=1s --timeout=1s --retries=5 \
	CMD curl \
	-X \
	POST \
	-s \
	-H 'Content-Type: application/json' \
	-d '{"jsonrpc":"2.0","id":"1","method":"net_listening","params":[]}' \
	http://127.0.0.1:8545

CMD ["anvil", "--base-fee", "0", "--block-time", "1", "--silent", "--host", "0.0.0.0", "--chain-id", "13370", "--load-state", "/usr/share/devnet/anvil_state.json"]
