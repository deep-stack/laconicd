#!/usr/bin/env bash
# Forwards all args to yarn on the sdk-test-runner container

if [ -n "$CERC_SCRIPT_DEBUG" ]; then
  set -x
fi

yarn_args=("--inspect-brk=8888")
yarn_args+=("${@:-test}")

# Get the key from laconic2d
laconic2d_key=$(
  yes | docker compose exec laconic2d laconic2d keys export alice --keyring-backend test --unarmored-hex --unsafe
)
# Set parameters for the test suite
cosmos_chain_id=laconic_9000-1
laconicd_rpc_endpoint=http://laconic2d:26657
laconicd_gql_endpoint=http://laconic2d:9473/api

docker compose exec laconic2d sh -c "curl --retry 10 --retry-delay 3 --retry-connrefused http://127.0.0.1:9473/api"

# Run tests
docker compose exec \
  -e COSMOS_CHAIN_ID="$cosmos_chain_id" \
  -e LACONICD_REST_ENDPOINT="$laconicd_rpc_endpoint" \
  -e LACONICD_GQL_ENDPOINT="$laconicd_gql_endpoint" \
  -e PRIVATE_KEY="$laconic2d_key" \
  sdk-test-runner yarn run "${yarn_args[@]}"
