# pggo
[![CI](https://github.com/bongnv/pggo/actions/workflows/ci.yml/badge.svg)](https://github.com/bongnv/pggo/actions/workflows/ci.yml) [![codecov](https://codecov.io/gh/bongnv/pggo/branch/main/graph/badge.svg?token=T9vvGtYE1d)](https://codecov.io/gh/bongnv/pggo) [![Go Report Card](https://goreportcard.com/badge/github.com/bongnv/pggo)](https://goreportcard.com/report/github.com/bongnv/pggo)

A code generator that allows you to write SQL queries in Go way.

## Development

- We use [pre-commit](https://pre-commit.com/) to format code & identify simple issues before submitting code to review:
  ```bash
  brew install pre-commit
  pre-commit install
  ```

- For unit tests, you can just run `go test` or `make test`.

- For integration tests, you'll need to setup docker:
  ```bash
  make setup-docker
  make test-integration
  ```

- Run `make help` to see all available commands with `make`.


## References

The project is inspired by:
- https://github.com/huandu/go-sqlbuilder
- https://github.com/lqs/sqlingo
- https://github.com/xo/xo
