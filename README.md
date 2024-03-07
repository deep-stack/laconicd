<div align="center">
  <h1> Laconic Network </h1>
</div>

![banner](docs/laconic.jpeg)

The Source of Proof. Laconic is a next generation data availability & verifiability layer with cryptographic proofs, powering internet-scale Web3 applications, built on Proof-of-Stake with fast-finality using the [Cosmos SDK](https://github.com/cosmos/cosmos-sdk/) which runs on top of [CometBFT](https://github.com/cometbft/cometbft) consensus engine.

## Installation

Install `laconic2d`:

  ```bash
  # install the laconic2d binary
  make install
  ```

## Usage

Run with a single node fixture:

  ```bash
  # start the chain
  ./scripts/init.sh

  # start the chain with data dir reset
  ./scripts/init.sh clean
  ```

## Tests

Run tests:

  ```bash
  # integration tests
  make test-integration

  # e2e tests
  make test-e2e
  ```
