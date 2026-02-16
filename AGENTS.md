# Repository Guidelines

## Project Structure

- `cmd/n8n/`: CLI entrypoint
- `internal/`: implementation packages
  - `cmd/`: Kong CLI commands
  - `api/`: REST client
  - `n8n/`: n8n API client and types
  - `secrets/`: Keyring-backed credential storage (API key + URL)
  - `outfmt/`: JSON/plain output formatting
  - `errfmt/`: User-friendly error formatting
  - `config/`: Platform-aware config paths

## Build, Test, and Development Commands

- `make` / `make build`: build `bin/n8n-cli`
- `make tools`: install pinned dev tools into `.tools/`
- `make fmt` / `make lint` / `make test` / `make ci`: format, lint, test, full local gate
- Hooks: `lefthook install` enables pre-commit checks

## Coding Style & Naming Conventions

- Formatting: `make fmt` (`goimports` local prefix `github.com/builtbyrobben/n8n-cli` + `gofumpt`)
- Output: keep stdout parseable (`--json` / `--plain`); send human hints/progress to stderr

## Testing Guidelines

- Unit tests: stdlib `testing` (and `httptest` where needed)
- All tests should use mocked HTTP servers (no live API calls in CI)

## Commit & Pull Request Guidelines

- Follow Conventional Commits + action-oriented subjects (e.g. `feat(cli): add --verbose flag`)
- Group related changes; avoid bundling unrelated refactors
- PRs should summarize scope, note testing performed, and mention user-facing changes

## Security

- Never commit API keys or credentials
- Use `--stdin` for credential input (avoid shell history leaks)
- Prefer OS keychain backends; use file backend only for headless environments
