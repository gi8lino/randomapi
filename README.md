# randomapi

**randomapi** is a tiny HTTP microservice that returns **one random element** from a JSON array.
It is designed to be extremely small, dependency-free, and production-ready.

Use it to serve:

- random **jokes**
- random **quotes**
- random **messages of the day**
- random **data items** of any shape (strings, objects, numbers…)

The service loads a JSON file at startup and exposes:

- `GET /random` → **one random element** from the JSON array
- `GET /healthz` → `"ok"` for liveness
- Optional `--route-prefix` support (e.g. `/api`)

> The sample `examples/jokes.json` is borrowed from the excellent
> [https://github.com/15Dkatz/official_joke_api](https://github.com/15Dkatz/official_joke_api)
> (MIT License).

## Configuration

`randomapi` is configured using **command-line flags** or **environment variables**
(prefixed with `RANDOMAPI_`).
Flags always take precedence over env vars.

| Flag               | Type   | Default          | Description                                                                 |
| ------------------ | ------ | ---------------- | --------------------------------------------------------------------------- |
| `--data-path`      | string | `/app/data.json` | Path to a JSON file containing a **JSON array** (any element type allowed). |
| `--listen-address` | string | `:8080`          | HTTP listen address for `/random` and `/healthz`.                           |
| `--route-prefix`   | string | _(empty)_        | Optional URL prefix to mount all endpoints under (e.g. `/api`).             |
| `--log-format`     | string | `text`           | Logging format: `text` or `json`.                                           |

### Environment Variables

Environment variables map directly to flags:

| Environment Variable       | Example               |
| -------------------------- | --------------------- |
| `RANDOMAPI_DATA_PATH`      | `/config/quotes.json` |
| `RANDOMAPI_LISTEN_ADDRESS` | `0.0.0.0:9090`        |
| `RANDOMAPI_ROUTE_PREFIX`   | `/randomapi`          |
| `RANDOMAPI_LOG_FORMAT`     | `json`                |

---

## Example Data File

Place a JSON array at the path provided to `--data-path`.

Example: `examples/jokes.json` (borrowed from 15Dkatz/official_joke_api):

```json
[
  {
    "type": "general",
    "setup": "Why don't scientists trust atoms?",
    "punchline": "Because they make up everything!"
  },
  {
    "type": "programming",
    "setup": "How many programmers does it take to change a light bulb?",
    "punchline": "None. It's a hardware problem."
  }
]
```

## API

### `GET /random`

Returns **one random element** from the loaded JSON file, exactly as it appears in the source.

```bash
curl http://localhost:8080/random
# → {"type":"general","setup":"...","punchline":"..."}
```

If `--route-prefix=/api` is used, the endpoint becomes:

```bash
GET /api/random
```

### `GET /healthz`

Simple liveness check:

```bash
curl http://localhost:8080/healthz
# → ok
```

---

## Run (local)

```bash
go run ./cmd/randomapi \
  --data-path=./examples/jokes.json \
  --listen-address=":8080"
```

Or using environment variables:

```bash
RANDOMAPI_DATA_PATH=./examples/jokes.json \
RANDOMAPI_LOG_FORMAT=json \
go run ./cmd/randomapi
```

---

## Kubernetes Deployment

You can find a complete example in the `examples/kubernetes` directory.

## Docker

```bash
docker run --rm \
  -p 8080:8080 \
  -v $(pwd)/examples/jokes.json:/data.json \
  ghcr.io/your-org/randomapi:latest \
  --data-path=/data.json
```

## License

This project is licensed under the Apache License.
The example jokes (`examples/jokes.json`) are borrowed from
[https://github.com/15Dkatz/official_joke_api](https://github.com/15Dkatz/official_joke_api) (MIT License).
