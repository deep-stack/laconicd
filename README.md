# laconic2d

Install `laconic2d`:

  ```bash
  # install the laconic2d binary
  make install
  ```

Run with a single node fixture:

  ```bash
  # start the chain
  ./scripts/init.sh

  # start the chain with data dir reset
  ./scripts/init.sh clean
  ```

Run tests:

  ```bash
  # integration tests
  make test-integration

  # e2e tests
  make test-e2e
  ```
