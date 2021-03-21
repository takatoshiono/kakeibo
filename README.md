[![backend CI](https://github.com/takatoshiono/kakeibo/workflows/backend%20CI/badge.svg)](https://github.com/takatoshiono/kakeibo/actions)
[![codecov](https://codecov.io/gh/takatoshiono/kakeibo/branch/main/graph/badge.svg?token=HSH2Wcy5C4)](https://codecov.io/gh/takatoshiono/kakeibo)

# Kakeibo

Kakeibo is a tool set to manage expenses.

## proto

We use [Buf](https://buf.build/) for Protobuf files.

### Installation

Buf:
- See [Installation](https://docs.buf.build/installation) page.

protoc-gen-go:
```
$ go install google.golang.org/protobuf/cmd/protoc-gen-go
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

### Generation
To generate stubs, run:
```
$ buf generate
```

### Lint
To run linting:
```
$ buf lint
```
