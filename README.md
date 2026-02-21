# dotld

CLI for domain availability and price search with Dynadot.

## What it does

- Searches domain availability and registration price.
- `dotld <domain.tld>` checks exact domains.
- `dotld <keyword>` (no TLD) checks mainstream suggestions:
  - `.com`, `.net`, `.org`, `.io`, `.ai`, `.co`, `.app`, `.dev`, `.sh`
- Prints compact tree output and supports `--json`.

## Example queries

```bash
# Single exact domain
dotld example.com

# Name without TLD (suggestions)
dotld acme

# Multiple exact domains
dotld example.com example.io

# JSON output
dotld example.com --json

# One-off key override
dotld example.com --dynadot-key your_dynadot_key
```

## Create Dynadot API Key

1. Log in to Dynadot.
2. Open `https://www.dynadot.com/account/domain/setting/api.html`.
3. Generate a production API key.
4. Export it in your shell:

```bash
export DYNADOT_API_PRODUCTION_KEY=your_key_here
```

## Setup

```bash
bun install
```

## Run

```bash
bun run dotld -- <domain.tld>
bun run dotld -- <keyword>
```

Direct key override:

```bash
bun run dotld -- <domain.tld> --dynadot-key <your_dynadot_key>
```

You can keep the export in your shell profile (`~/.zshrc`, `~/.bashrc`, or fish config) for persistence.

JSON mode:

```bash
bun run dotld -- <domain.tld> --json
```

## Local install

```bash
bun run link:local
dotld <domain.tld>
```

## curl | bash installer

Build release binaries:

```bash
bun run build:release
```

Automated release (runs lint/check/test, builds binaries, updates `CHANGELOG.md`, tags, pushes, and creates GitHub release assets):

```bash
bun run release
```

Optional release modes:

```bash
bun run release --dry-run
bun run release --bump minor
bun run release --version 1.0.0
bun run release --notes "custom release notes"
```

Install from GitHub Releases:

```bash
DOTLD_REPO=your-org/dotld curl -fsSL https://raw.githubusercontent.com/your-org/dotld/main/scripts/install.sh | bash
```

Install and run immediately:

```bash
DOTLD_REPO=your-org/dotld curl -fsSL https://raw.githubusercontent.com/your-org/dotld/main/scripts/install.sh | bash -s -- -- <domain.tld>
```

## Validation

```bash
bun lint
bun check
bun test
```

## Dynadot limits note

Regular Dynadot accounts are limited to 1 domain per `search` command and 60 requests/min.
This CLI intentionally uses single-domain requests for predictable behavior.
