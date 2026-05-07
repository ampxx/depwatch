# depwatch

A daemon that monitors dependency version changes in Go modules and sends alerts via webhooks.

---

## Installation

```bash
go install github.com/yourname/depwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/depwatch.git && cd depwatch && go build -o depwatch .
```

---

## Usage

Create a configuration file `depwatch.yaml`:

```yaml
module_path: /path/to/your/go/project
interval: 6h
webhooks:
  - url: https://hooks.slack.com/services/your/webhook/url
    on: version_change
```

Start the daemon:

```bash
depwatch --config depwatch.yaml
```

depwatch will poll your `go.mod` and `go.sum` files at the configured interval and fire a webhook payload whenever a dependency version change is detected.

### Example Webhook Payload

```json
{
  "module": "github.com/gin-gonic/gin",
  "previous_version": "v1.9.0",
  "current_version": "v1.9.1",
  "detected_at": "2024-11-03T14:32:00Z"
}
```

---

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `depwatch.yaml` | Path to config file |
| `--once` | `false` | Run a single check and exit |
| `--log-level` | `info` | Log verbosity (`debug`, `info`, `warn`) |

---

## License

MIT © 2024 yourname