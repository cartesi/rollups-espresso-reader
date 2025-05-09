# Cartesi Rollups Espresso Reader

Input reader implementation for the [Cartesi Rollups Node](https://github.com/cartesi/rollups-node) that follows an [Espresso Network](https://docs.espressosys.com/network) to fetch input data.

## Introduction

[Espresso Systems](https://www.espressosys.com/) provides a decentralized sequencer and data availability system, which can be very useful to scale layer-2 rollup solutions.
[Cartesi](https://cartesi.io) provides an [app-specific rollups solution](https://docs.cartesi.io/cartesi-rollups/) which could particularly benefit from both.

The [Cartesi Rollups Node](https://github.com/cartesi/rollups-node) is the component responsible for fetching inputs for Cartesi applications and processing them.
The Node supports a configuration to turn off its input reading functionality, allowing an alternative input reader implementation to be used instead.

In this context, this repository contains an input reader implementation that pulls data from an Espresso Network while retaining the ability to fetch inputs from the Cartesi Rollups InputBox contract on the base layer (as they usually do).
As such, applications can take advantage of Espresso's higher throughput and lower fees and latency, while still being able to interact with the base layer (e.g., to manage L1 assets).

## Architecture

Check the [architecture](docs/architecture.md) page for more details about how the integration of Espresso with Cartesi Rollups works.

## App development

Check the [development](docs/development.md) page for details about how to develop Cartesi applications using Espresso.

## Running locally

The Espresso Reader is intended to be executed alongside a Cartesi Rollups Node.

In order to run it locally, it is necessary to instantiate an environment with all the necessary components: a local blockchain network, a local Espresso network, and a Cartesi Rollups Node alongside its database and the Espresso Reader itself.

To build all necessary components:

```bash
docker compose build
```

Then, you can execute everything by running:

```bash
docker compose up -d
```

Logs for the Cartesi Node with the Espresso Reader can be tracked by typing:

```bash
docker compose logs cartesi_node_espresso -f
```

A sample Echo application can be deployed on your local node by executing:

```bash
INPUT_BOX_ADDRESS="0xB6b39Fb3dD926A9e3FBc7A129540eEbeA3016a6c" \
ESPRESSO_STARTING_BLOCK="0" \
ESPRESSO_NAMESPACE="55555" \
DATA_AVAILABILITY=$(cast calldata \
    "InputBoxAndEspresso(address,uint256,uint32)" \
    $INPUT_BOX_ADDRESS $ESPRESSO_STARTING_BLOCK $ESPRESSO_NAMESPACE); \
docker compose exec cartesi_node_espresso cartesi-rollups-cli app deploy -v \
    --salt 0000000000000000000000000000000000000000000000000000000000000000 \
    --data-availability $DATA_AVAILABILITY \
    -n echo-dapp \
    -t applications/echo-dapp/
```

Once deployed, an L1 InputBox input can be sent using cast:

```bash
INPUT=0xdeadbeef; \
INPUT_BOX_ADDRESS=0xB6b39Fb3dD926A9e3FBc7A129540eEbeA3016a6c; \
APPLICATION_ADDRESS=0x01e800bbE852aeb27cE65604709134Ea63782c6B; \
cast send \
    --mnemonic "test test test test test test test test test test test junk" \
    --rpc-url "http://localhost:8545" \
    $INPUT_BOX_ADDRESS "addInput(address,bytes)(bytes32)" $APPLICATION_ADDRESS $INPUT
```

### Running with an app deployed on Sepolia testnet

You can also run your local node to process apps deployed on Sepolia testnet.
In this scenario, the Espresso network to be used should be [Decaf](https://docs.espressosys.com/network/releases/testnets/decaf-testnet).

The file [env.nodev2-sepolia-decaf](./ci/env.nodev2-sepolia-decaf) contains the basic environment variable settings in order to setup your local Node + Espresso Reader to use Sepolia + Decaf.
Make sure to edit the first lines of that file to specify an appropriate Sepolia account and blockchain gateway.

To start up your node, execute:

```bash
NODE_ENV_FILE=./ci/env.nodev2-sepolia-decaf docker compose up -d db cartesi_node_espresso
```

You should then deploy a testnet application (similarly to sample Echo deployment described above), or register an existing testnet application on your node.
To register, execute the following:

```bash
eval $(cat ./ci/env.nodev2-sepolia-decaf)
docker compose exec cartesi_node_espresso cartesi-rollups-cli app register -v \
    -c 0xConsensusContractAddress \
    -a 0xApplicationAddress \
    -n app-name \
    -t path/to/app/template/
```



## Building and testing with a local Cartesi Node repository

First build the Espresso Reader itself:

```bash
go build
```

To run it alongisde a Cartesi Rollups Node checked out from its [repository](https://github.com/cartesi/rollups-node/releases/tag/v2.0.0-alpha.4):

```bash
cd <path-to-cartesi-rollups-node>
make
make run-postgres && make migrate
./cartesi-rollups-node
```

Then, run the Espresso Reader to read inputs from Espresso and write them to the Node's database (make sure to adjust environment variables as appropriate to match the Node's configuration):

```bash
eval $(make env)
make migrate
./rollups-espresso-reader
```

Alternatively, to run automated integration tests:

```bash
go test github.com/cartesi/rollups-espresso-reader/internal/espressoreader -v
```
