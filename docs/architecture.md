# Cartesi+Espresso Architecture

This section describes the concepts and architecture of the integration of Espresso with Cartesi Rollups.

## Concepts

The Cartesi+Espresso integration is based on the concept that inputs to Cartesi applications are of two fundamentally different natures:

### L2 transactions

These refer to common interactions of users with the application, and refer to application-specific actions such as “attack goblin”, “swap token”, “post message”, etc.; these transactions do not require any direct information or logic from the base layer;

### L1->L2 messages

These refer to information that is relayed from the base layer to the rollup application, such as informing about deposits done via the Portals, relaying the app’s address, ensuring base layer validation for a given input, etc.

## Architecture

TODO

