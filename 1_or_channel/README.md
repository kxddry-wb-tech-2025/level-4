## or-channel (Go)

Minimal implementation and example of the "or-channel" pattern in Go. It combines multiple `<-chan struct{}` into a single channel that closes when any of the inputs close.

### Layout
- `pkg/or`: library code and tests
- `cmd/example`: small runnable example

### Run the example
```bash
go run ./cmd/example
```

### Run tests
```bash
go test ./...
```

### Go version
Defined in `go.mod`.


