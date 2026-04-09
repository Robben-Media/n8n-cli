# n8n-cli

n8n CLI — Workflow automation from the command line.

## Installation

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/Robben-Media/n8n-cli/releases).

### Build from Source

```bash
git clone https://github.com/Robben-Media/n8n-cli.git
cd n8n-cli
go build ./cmd/n8n
```

## Configuration

n8n-cli requires an API key and your n8n instance URL. Credentials are stored securely in your system keyring.

**Store credentials:**

```bash
n8n-cli auth set-key --url https://n8n.example.com
```

The CLI will prompt for your API key interactively. You can also pipe it:

```bash
echo "your-api-key" | n8n-cli auth set-key --url https://n8n.example.com
```

**Environment variable override:**

```bash
export N8N_API_KEY="your-api-key"
export N8N_URL="https://n8n.example.com"
```

**Check status:**

```bash
n8n-cli auth status
```

**Remove credentials:**

```bash
n8n-cli auth remove
```

## Commands

### auth

Manage API credentials.

| Command | Description |
|---------|-------------|
| `auth set-key --url <url>` | Store API key and n8n instance URL in keyring |
| `auth status` | Show authentication status |
| `auth remove` | Remove stored credentials |

### workflows

| Command | Description |
|---------|-------------|
| `workflows list` | List workflows |
| `workflows get <id>` | Get a workflow by ID |
| `workflows activate <id>` | Activate a workflow |
| `workflows deactivate <id>` | Deactivate a workflow |
| `workflows delete <id>` | Delete a workflow |

**Flags (list):** `--active` (filter by active status), `--tags` (filter by tag name), `--limit` (default 20), `--cursor` (pagination cursor)

### executions

| Command | Description |
|---------|-------------|
| `executions list` | List executions |
| `executions get <id>` | Get an execution by ID |
| `executions delete <id>` | Delete an execution |
| `executions retry <id>` | Retry a failed execution |

**Flags (list):** `--workflow-id`, `--status` (success, error, waiting), `--limit` (default 20), `--cursor` (pagination cursor)

### credentials

| Command | Description |
|---------|-------------|
| `credentials list` | List credentials |

### tags

| Command | Description |
|---------|-------------|
| `tags list` | List tags |
| `tags create <name>` | Create a tag |
| `tags delete <id>` | Delete a tag |

### variables

| Command | Description |
|---------|-------------|
| `variables list` | List variables |
| `variables create --key <k> --value <v>` | Create a variable |
| `variables delete <id>` | Delete a variable |

### webhooks

| Command | Description |
|---------|-------------|
| `webhooks trigger --path <path>` | Trigger a webhook |

**Flags:** `--method` (default POST), `--data` (JSON data to send)

### health

| Command | Description |
|---------|-------------|
| `health` | Check n8n instance health |

## Global Flags

| Flag | Description |
|------|-------------|
| `--json` | Output JSON to stdout (best for scripting) |
| `--plain` | Output stable, parseable text to stdout (TSV; no colors) |
| `--verbose` | Enable verbose logging |
| `--force` | Skip confirmations for destructive commands |
| `--no-input` | Never prompt; fail instead (useful for CI) |
| `--color` | Color output: auto, always, or never (default auto) |

## License

MIT
