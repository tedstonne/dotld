# dotld

`dotld` is a CLI for domain availability and registration price search through [Dynadot](https://www.dynadot.com/account/domain/setting/api.html).

## Features

- Exact lookup: `dotld example.com`
- No-TLD suggestions: `dotld keyword`
- Compact terminal output and `--json` mode
- User-owned key flow (`DYNADOT_API_PRODUCTION_KEY` or `--dynadot-key`)
- Auto-persists API key to `~/.config/dotld/config.json`

## Requirements

- Dynadot production API key

## Quick Start

1. Generate your key in Dynadot:
   - https://www.dynadot.com/account/domain/setting/api.html
2. Run with your key (auto-saved for next time):

```bash
dotld --dynadot-key your_key_here example.com
```

3. Subsequent runs need no key:

```bash
dotld example.com
```

Or export it in your shell:

```bash
export DYNADOT_API_PRODUCTION_KEY=your_key_here
dotld example.com
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

## Claude Code Skill

Install as a Claude Code skill to let Claude search domains for you:

```bash
# skill.sh
npx skills add tedstonne/dotld

# ClawHub
clawhub install dotld
```

## Install

From GitHub Releases:

```bash
curl -fsSL https://raw.githubusercontent.com/tedstonne/dotld/main/scripts/install.sh | bash
```

Install and run immediately:

```bash
curl -fsSL https://raw.githubusercontent.com/tedstonne/dotld/main/scripts/install.sh | bash -s -- -- example.com
```

## Build from Source

```bash
go build -o dotld .
```

Cross-platform binaries:

```bash
just build
```

## Release (Maintainers)

```bash
just release          # patch (default)
just release minor
just release major
```

## Validation

```bash
just test
just lint
```

## Dynadot Limits

Regular Dynadot accounts are limited to 1 domain per `search` command and 60 requests/min.
`dotld` uses single-domain requests for predictable behavior.
