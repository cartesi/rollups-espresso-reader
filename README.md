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

## App Development

Check the [development](docs/development.md) page for  details about how to develop Cartesi applications using Espresso.

## Building

To build:

```bash
go build
```

To run:

```bash
eval $(make env)
make migrate
./espresso-reader
```

To run automated integrated tests, 
```
(cartesi-rollups-node)
make run-postgres && make migrate
./cartesi-rollups-node

(rollups-espresso-reader)
make migrate
go test github.com/cartesi/rollups-espresso-reader/internal/espressoreader -v
```
