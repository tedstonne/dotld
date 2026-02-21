# Domain Search CLI + Affiliate API Project Plan

## Goal
Build a zero-config CLI for users to search domain availability + pricing across TLDs, without requiring users to provide registrar API keys. Monetize via affiliate links.

## Product Summary
- User runs CLI locally (`domain search ...`)
- CLI calls your hosted API
- Hosted API queries registrar APIs using your private credentials
- API returns normalized results: availability, price, currency, buy link
- CLI prints clean table and JSON output option

## Recommended Stack
- **Backend API**: Cloudflare Workers (TypeScript)
- **CLI**: Bun + TypeScript
- **Cache/State**: Cloudflare KV (short TTL), optional Workers Cache API
- **Registrar v1**: Dynadot (simplest search+price integration)
- **Future adapters**: GoDaddy, Namecheap

## Why This Architecture
- Users need no API keys
- Your keys remain server-side and secure
- Easy to add abuse controls and caching
- One language across API + CLI for faster iteration

---

## Scope v1 (MVP)

### In Scope
- Dynadot integration for domain search + price
- CLI command: `domain search <domains...>`
- API endpoint: `POST /v1/search`
- Affiliate link generation for available domains
- Output formats: table + JSON
- Basic auth/rate limiting/caching

### Out of Scope (v1)
- Domain purchasing from CLI
- Account management
- Multi-registrar aggregation
- Fancy analytics dashboard

---

## High-Level Architecture

1. CLI collects domains and sends request to Worker API
2. Worker validates request and applies rate limiting
3. Worker checks KV/cache for recent results
4. Worker queries Dynadot API for cache misses
5. Worker normalizes response + adds affiliate buy URLs
6. Worker returns structured JSON
7. CLI renders table (or JSON with `--json`)

---

## API Contract (v1)

### Endpoint
`POST /v1/search`

### Request
```json
{
  "domains": ["example.com", "brand.ai"],
  "currency": "USD"
}
```

### Response
```json
{
  "results": [
    {
      "domain": "example.com",
      "available": false,
      "price": null,
      "currency": "USD",
      "buyUrl": null,
      "source": "dynadot",
      "cached": true
    },
    {
      "domain": "brand.ai",
      "available": true,
      "price": "79.99",
      "currency": "USD",
      "buyUrl": "https://<affiliate-link>",
      "source": "dynadot",
      "cached": false
    }
  ]
}
```

### Validation Rules
- Max domains/request: 20 (initial)
- Deduplicate input
- Normalize lowercase
- Domain format validation (reject invalids)
- Allowed currencies list (start with `USD`)

---

## CLI UX (v1)

### Commands
- `domain search example.com brand.ai`
- `domain search --file domains.txt`
- `domain search example.com --json`

### Flags
- `--json` machine-readable output
- `--currency USD`
- `--api https://api.yourdomain.com`
- `--timeout 10s`

### Output Columns
- Domain
- Available (yes/no)
- Price
- Currency
- Registrar
- Buy URL (affiliate)

---

## Security Plan

### Secrets
- Store `DYNADOT_API_KEY` in Cloudflare Workers secrets
- Never expose registrar key to CLI or logs

### API Access Control
- Option A: Public endpoint + strict per-IP limits
- Option B: CLI client token (`X-Client-Key`) + per-token limits
- Recommendation: **Option B** for better abuse control

### Abuse Prevention
- Rate limit by IP + token
- Request body size limit
- Max domains/request
- Reject high-frequency duplicate requests

---

## Caching + Cost Control
- Cache key: `<domain>:<currency>`
- TTL: 60 seconds (tune later)
- Cache misses only call registrar API
- Batch lookups in single registrar call when possible
- Track cache hit rate in logs

---

## Affiliate + Compliance
- Include affiliate disclosure in CLI help/README
- Ensure link format follows affiliate program terms
- Avoid misleading branding ("official registrar CLI" claims)
- Add terms page for your API usage policy

---

## Observability
- Log structured events:
  - request id
  - domains count
  - cache hit/miss
  - registrar latency
  - error type
- Track metrics:
  - requests/day
  - cost per 1k queries
  - top TLDs
  - affiliate click-outs (if tracked)

---

## Milestones

### Milestone 1: Foundation (Day 1-2)
- Set up Worker project + route
- Add secret management
- Build `POST /v1/search` skeleton
- Add request validation

### Milestone 2: Dynadot Integration (Day 3)
- Implement Dynadot search call
- Parse and normalize response
- Add affiliate URL generation

### Milestone 3: CLI MVP (Day 4)
- Build `domain search` command
- Add table + JSON output
- Support multiple domains

### Milestone 4: Hardening (Day 5)
- Add rate limiting
- Add KV caching
- Improve error handling + retries
- Write docs and usage examples

### Milestone 5: Launch Readiness (Day 6-7)
- Smoke tests
- Basic load test
- Publish CLI package/binary
- Create changelog + v0.1.0 release

---

## Testing Plan

### Backend Tests
- Input validation tests
- Dynadot response parser tests
- Error mapping tests
- Cache behavior tests
- Rate limit tests

### CLI Tests
- Command parsing tests
- JSON output snapshot tests
- Table output smoke tests
- API timeout/error handling tests

### E2E
- Search known domains
- Mixed available/unavailable domains
- Invalid domain list handling

---

## Error Handling Strategy
- Registrar timeout -> return partial with per-domain error
- Invalid domain -> return structured validation error
- Rate limited -> 429 with retry hint
- Upstream failure -> 502 with safe message

---

## Future Roadmap (Post-MVP)
1. Add GoDaddy adapter
2. Add Namecheap adapter
3. Add registrar preference flag (`--registrar`)
4. Add best-price aggregation across registrars
5. Add TLD suggestion engine
6. Add click-tracking redirect endpoint

---

## Risks and Mitigations

### Risk: Affiliate attribution loss
- Mitigation: always return direct affiliate buy link; test conversion path

### Risk: API abuse
- Mitigation: auth tokens + IP rate limiting + quotas + caching

### Risk: Registrar API changes
- Mitigation: adapter layer + integration tests + feature flags

### Risk: Pricing mismatch at checkout
- Mitigation: label as "latest API quote"; include timestamp

---

## Definition of Done (v1)
- Users can install CLI and run search without providing registrar keys
- API returns availability + price + affiliate buy URL
- Basic rate limits and cache are active
- README includes setup and disclosure
- MVP deployed on Cloudflare Workers with monitoring

---

## Setup Checklist (Pre-Implementation)
- [ ] Create Dynadot production API key
- [ ] Confirm affiliate link format and terms
- [ ] Create Cloudflare account/project
- [ ] Provision Worker + KV namespace
- [ ] Decide API auth mode (public vs client-key)
- [ ] Decide CLI distribution (bun package vs binary)
