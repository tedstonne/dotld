# 🌐 dotld

Domain TLD search for your CLI and your AI agent. Availability and prices in one command, or one skill install away from your agent.

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
├─ acme.sh  · Taken
└─ acme.so  · Taken
```

Type a keyword, get every TLD that matters. Available domains show price and a link to buy.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/tedstonne/dotld/main/scripts/install.sh | bash
```

Requires a [Dynadot production API key](https://www.dynadot.com/account/domain/setting/api.html) (free with any account). Your first run saves the key locally. No config files to manage.

```bash
dotld --dynadot-key YOUR_KEY acme
```

Every run after that is just `dotld <keyword>`.

## Usage

Check a specific domain:

```text
$ dotld acme.xyz
acme.xyz · $9.99 · https://www.dynadot.com/domain/search?domain=acme.xyz
```

Search across TLDs with a keyword:

```text
$ dotld startup
```

Get structured output for scripts and pipelines:

```text
$ dotld bigpickle.com --json
```

Override your saved key for a one-off lookup:

```text
$ dotld bigpickle.com --dynadot-key OTHER_KEY
```

## Skills

Give your AI agent the ability to search domains mid-conversation.

```bash
npx skills add tedstonne/dotld
```

Works with [skills.sh](https://skills.sh) compatible agents: Claude Code, OpenCode, Codex, Gemini CLI, and others.

> **Note:** The skill uses your local `dotld` install. A Dynadot production API key must be configured on your machine.