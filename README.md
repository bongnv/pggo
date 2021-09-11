# pggo
A code generator to write SQL queries in Go.

## Development

- We use [pre-commit](https://pre-commit.com/) to format code & identify simple issues before submitting code to review:
  ```bash
  brew install pre-commit
  pre-commit install
  ```

## Testing

### Integration

- Before running integration tests, you'll need to setup docker
  ```bash
  make setup-docker
  ```

- Run integration tests
  ```bash
  make test-integration
  ```
