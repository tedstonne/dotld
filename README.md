# dotld

CLI for fast Dynadot availability checks with clean terminal output.

## What it does

- `dotld name.tld` checks a single domain.
- `dotld name` (no TLD) checks mainstream suggestions:
  - `.com`, `.net`, `.org`, `.io`, `.ai`, `.co`, `.app`, `.dev`, `.sh`
- Prints compact tree output and supports `--json`.

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
bun run dotld -- murk.ink
bun run dotld -- murk
```

Direct key override:

```bash
bun run dotld -- murk --dynadot-key your_key_here
```

You can keep the export in your shell profile (`~/.zshrc`, `~/.bashrc`, or fish config) for persistence.

JSON mode:

```bash
bun run dotld -- murk --json
```

## Local install

```bash
bun run link:local
dotld murk
```

## curl | bash installer

Build release binaries:

```bash
bun run build:release
```

Install from GitHub Releases:

```bash
DOTLD_REPO=your-org/dotld curl -fsSL https://raw.githubusercontent.com/your-org/dotld/main/scripts/install.sh | bash
```

Install and run immediately:

```bash
DOTLD_REPO=your-org/dotld curl -fsSL https://raw.githubusercontent.com/your-org/dotld/main/scripts/install.sh | bash -s -- -- murk
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
