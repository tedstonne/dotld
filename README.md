# dotld

`dotld` is a CLI for domain availability and registration price search through [Dynadot](https://www.dynadot.com/account/domain/setting/api.html).

## Features

- Exact lookup: `dotld example.com`
- No-TLD suggestions: `dotld keyword`
- Compact terminal output and `--json` mode
- User-owned key flow (`DYNADOT_API_PRODUCTION_KEY` or `--dynadot-key`)

## Requirements

- Bun
- Dynadot production API key

## Quick Start

1. Generate your key in Dynadot:
   - https://www.dynadot.com/account/domain/setting/api.html
2. Export it in your shell:

```bash
export DYNADOT_API_PRODUCTION_KEY=your_key_here
```

3. Install dependencies and run:

```bash
bun install
bun run dotld -- example.com
```

## Examples

Exact domain:

```text
$ dotld example.com
example.com · Taken
```

Keyword suggestions (no TLD input):

```text
$ dotld acme
acme
├─ acme.com · Taken
├─ acme.net · Taken
├─ acme.org · Taken
├─ acme.io  · $39.99 · https://www.dynadot.com/domain/search?domain=acme.io
├─ acme.ai  · Taken
├─ acme.co  · Taken
├─ acme.app · Taken
├─ acme.dev · Taken
└─ acme.sh  · Taken
```

JSON mode:

```text
$ dotld example.com --json
{
  "results": [
    {
      "domain": "example.com",
      "available": false,
      "price": null,
      "currency": "USD",
      "buyUrl": null,
      "source": "dynadot",
      "cached": false,
      "quotedAt": "2026-02-21T00:00:00.000Z"
    }
  ]
}
```

One-off key override:

```text
$ dotld example.com --dynadot-key your_dynadot_key
example.com · Taken
```

## Local Install

```bash
bun run link:local
dotld example.com
```

## Installer

Install from GitHub Releases:

```bash
DOTLD_REPO=your-org/dotld curl -fsSL https://raw.githubusercontent.com/your-org/dotld/main/scripts/install.sh | bash
```

Install and run immediately:

```bash
DOTLD_REPO=your-org/dotld curl -fsSL https://raw.githubusercontent.com/your-org/dotld/main/scripts/install.sh | bash -s -- -- example.com
```

## Release (Maintainers)

Build binaries only:

```bash
bun run build:release
```

Automated release (lint, check, test, build, changelog, tag, push, GitHub release assets):

```bash
bun run release
```

Other release options:

```bash
bun run release --dry-run
bun run release --bump minor
bun run release --version 1.0.0
bun run release --notes "custom release notes"
```

## Validation

```bash
bun lint
bun check
bun test
```

## Dynadot Limits

Regular Dynadot accounts are limited to 1 domain per `search` command and 60 requests/min.
`dotld` uses single-domain requests for predictable behavior.
