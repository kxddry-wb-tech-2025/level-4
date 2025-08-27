## Distributed Grep (client + server)

A minimal distributed "grep". The client splits input into chunks, sends them to multiple server instances, and merges the results to match the behavior of system `grep`.

- `server/`: HTTP service handling grep tasks
- `client/`: CLI that orchestrates distributed grep

### Prerequisites
- Go 1.25 ('cause of go.mod)
- Docker and Docker Compose (optional)

## Quick start with Docker Compose
Bring up three servers and optionally run a demo client in containers.

### Start servers
```bash
docker compose up -d server1 server2 server3
```
Servers will listen on:
- `server1:8081`
- `server2:8082`
- `server3:8083`

### Verify health
```bash
curl -i http://localhost:8081/health
```

### Run a one-off demo client (profile)
```bash
# Creates a small input file and runs the client against all servers (-n for line numbers)
docker compose --profile demo up client
```

### View logs
```bash
docker compose logs -f
```

### Stop everything
```bash
docker compose down
```

The `docker-compose.yml` uses the official `golang:1.25` image and runs modules with `go run`, so no Dockerfiles are needed.

## Running locally (no Docker)
Open two terminals: one for servers, one for the client.

### Start servers
```bash
# Terminal A
cd server

# build the server
go build -o server ./cmd/app

# Start three servers on ports 8080..8082
./server -d 
./server -d
./server -d
```
Health check endpoint:
```bash
curl -i http://localhost:8080/health
```

### Run client
```bash
# Terminal B
cd client

# Basic usage: grep 'foo' in file.txt using 3 servers
go run ./cmd/app --addrs 127.0.0.1:8080,127.0.0.1:8081,127.0.0.1:8082 foo file.txt
```

The client supports stdin when no files are provided:
```bash
echo -e "alpha\nbeta\nfoo" | go run ./cmd/app --addrs 127.0.0.1:8080,127.0.0.1:8081,127.0.0.1:8082 foo
```

### Building binaries
```bash
# Server
cd server
go build -o bin/server ./cmd/app
./bin/server -port=8081

# Client
cd client
go build -o bin/client ./cmd/app
./bin/client --addrs 127.0.0.1:8081 foo file.txt
```

## CLI flags (client)
The client mirrors a subset of `grep` flags and adds distributed options:

- **-A, --after NUM**: Print NUM lines of trailing context
- **-B, --before NUM**: Print NUM lines of leading context
- **-C, --context NUM**: Set both before and after context to NUM
- **-v, --invert**: Invert match selection
- **-i, --ignore-case**: Case-insensitive match
- **-c, --count**: Print only count of selected lines per file
- **-F, --fixed-string**: PATTERN is a literal string, not regex
- **-n, --print-numbers**: Print line numbers
- **--addrs host:port[,host:port...]**: Comma-separated server addresses (required)
- **--quorum N**: Minimum successful server responses (default: majority)

Examples:
```bash
# Case-insensitive match with context and line numbers
./client -i -C 2 -n --addrs 127.0.0.1:8081,127.0.0.1:8082 foo file.txt

# Count only
./client -c --addrs 127.0.0.1:8081,127.0.0.1:8082 foo file.txt
```

## Server endpoints
- `POST /grep` — accepts a task containing lines and returns found blocks
- `GET /health` — returns 204 when ready

## Integration tests
Integration tests compare the distributed client output against system `grep` across multiple scenarios.

Run tests from the client module:
```bash
cd client
go test -v
```

## How it works (brief)
- The client probes the `--addrs` for health to determine alive servers.
- The input is split into chunks; each chunk is sent with necessary context.
- Servers perform local matching (regex or fixed string) and return matching blocks.
- The client merges blocks and prints in file-order; with `-c`, it aggregates counts from all servers.
- Quorum defaults to a simple majority of alive servers.

## Troubleshooting
- **Ports busy**: change `-port` values or stop existing processes.
- **No output**: verify `--addrs` are correct and `/health` returns 204.
- **Different results from grep**: only a subset of `grep` is implemented.
- **Docker networking**: use `localhost:PORT` from host, or `server1:8081` inside Compose network.


