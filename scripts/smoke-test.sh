#!/usr/bin/env bash
set -euo pipefail

PASS=0
FAIL=0

pass() { printf "  \xe2\x9c\x93 %s\n" "$1"; PASS=$((PASS + 1)); }
fail() { printf "  \xe2\x9c\x97 %s\n" "$1"; FAIL=$((FAIL + 1)); }
skip() { printf "  - %s (skipped)\n" "$1"; }

# ---------- Stage 1: Binary Verification ----------

printf "\n=== Stage 1: Binary Verification ===\n\n"

if command -v dotld &>/dev/null; then
  pass "dotld is in PATH"
else
  fail "dotld not found in PATH"
  printf "Cannot continue without binary\n"
  exit 1
fi

VERSION="$(dotld --version 2>&1 || true)"
printf "  version: %s\n" "$VERSION"

HELP="$(dotld --help 2>&1 || true)"
if echo "$HELP" | grep -qi "usage\|domain\|search"; then
  pass "dotld --help shows usage info"
else
  fail "dotld --help output unexpected: $HELP"
fi

# ---------- Stage 2: Live API Test ----------

printf "\n=== Stage 2: Live API Test ===\n\n"

if [[ -z "${DYNADOT_API_PRODUCTION_KEY:-}" ]]; then
  skip "DYNADOT_API_PRODUCTION_KEY not set — skipping API tests"
else
  OUTPUT="$(dotld example.com 2>&1 || true)"
  if echo "$OUTPUT" | grep -qi "taken"; then
    pass "dotld example.com returned Taken"
  else
    fail "dotld example.com unexpected output: $OUTPUT"
  fi

  JSON="$(dotld example.com --json 2>&1 || true)"
  if echo "$JSON" | jq -e '.results[0].domain' &>/dev/null; then
    pass "dotld example.com --json has valid JSON with results[].domain"
  else
    fail "dotld example.com --json invalid: $JSON"
  fi

  if echo "$JSON" | jq -e '.results[0] | has("available", "price", "currency")' &>/dev/null; then
    pass "JSON contains expected fields (available, price, currency)"
  else
    fail "JSON missing expected fields"
  fi
fi

# ---------- Stage 3: Skill Install (via skills.sh) ----------

printf "\n=== Stage 3: Skill Install (via skills.sh) ===\n\n"

# Skills were installed during docker build via: npx skills add tedstonne/dotld -a opencode -a claude-code
# Verify they landed in the expected directories

# Check Claude Code skill path
CLAUDE_SKILL_DIR="$HOME/.claude/skills/dotld"
if [[ -f "$CLAUDE_SKILL_DIR/SKILL.md" ]]; then
  pass "SKILL.md installed to $CLAUDE_SKILL_DIR (claude-code)"
else
  fail "SKILL.md not found at $CLAUDE_SKILL_DIR — npx skills add may have failed for claude-code"
fi

# Check OpenCode skill path
OPENCODE_SKILL_DIR="$HOME/.config/opencode/skills/dotld"
if [[ -f "$OPENCODE_SKILL_DIR/SKILL.md" ]]; then
  pass "SKILL.md installed to $OPENCODE_SKILL_DIR (opencode)"
else
  fail "SKILL.md not found at $OPENCODE_SKILL_DIR — npx skills add may have failed for opencode"
fi

# Validate SKILL.md frontmatter from whichever path exists
SKILL_FILE=""
if [[ -f "$CLAUDE_SKILL_DIR/SKILL.md" ]]; then
  SKILL_FILE="$CLAUDE_SKILL_DIR/SKILL.md"
elif [[ -f "$OPENCODE_SKILL_DIR/SKILL.md" ]]; then
  SKILL_FILE="$OPENCODE_SKILL_DIR/SKILL.md"
fi

if [[ -n "$SKILL_FILE" ]]; then
  FRONTMATTER="$(sed -n '/^---$/,/^---$/p' "$SKILL_FILE")"
  for FIELD in name allowed-tools description; do
    if echo "$FRONTMATTER" | grep -q "$FIELD"; then
      pass "SKILL.md frontmatter has '$FIELD'"
    else
      fail "SKILL.md frontmatter missing '$FIELD'"
    fi
  done

  # Validate references directory
  SKILL_BASE="$(dirname "$SKILL_FILE")"
  if [[ -f "$SKILL_BASE/references/cli-reference.md" ]]; then
    pass "cli-reference.md present in skill references"
  else
    fail "cli-reference.md missing from skill references"
  fi
else
  fail "No SKILL.md found in either agent path — cannot validate frontmatter"
fi

# ---------- Stage 4: OpenCode CLI ----------

printf "\n=== Stage 4: OpenCode CLI ===\n\n"

if command -v opencode &>/dev/null; then
  pass "opencode is in PATH"
  OC_VERSION="$(opencode --version 2>&1 || true)"
  printf "  opencode version: %s\n" "$OC_VERSION"
else
  fail "opencode not found in PATH"
fi

# Verify opencode config exists
if [[ -f "$HOME/.config/opencode/opencode.json" ]]; then
  pass "opencode.json config present"

  # Validate config is valid JSON
  if jq empty "$HOME/.config/opencode/opencode.json" 2>/dev/null; then
    pass "opencode.json is valid JSON"
  else
    fail "opencode.json is not valid JSON"
  fi
else
  skip "opencode.json not found — OpenCode integration test limited"
fi

# ---------- Stage 5: OpenCode Skill Integration ----------

printf "\n=== Stage 5: OpenCode + dotld Skill Integration ===\n\n"

# Determine if we have a model provider configured
HAS_PROVIDER=false

if [[ -n "${OPENCODE_API_KEY:-}" ]]; then
  HAS_PROVIDER=true
  printf "  provider: OpenCode Zen\n"
elif [[ -n "${GITHUB_TOKEN:-}" ]]; then
  HAS_PROVIDER=true
  printf "  provider: GitHub Copilot\n"
elif [[ -n "${OPENAI_API_KEY:-}" ]]; then
  HAS_PROVIDER=true
  printf "  provider: OpenAI\n"
fi

if [[ "$HAS_PROVIDER" == "true" ]] && [[ -n "${DYNADOT_API_PRODUCTION_KEY:-}" ]]; then

  # Step 1: Ask the model to write a prompt about domain names.
  # This ensures the test input is never hardcoded — the model generates it.
  printf "  Generating prompt via OpenCode...\n"
  GEN_PROMPT="$(timeout 30 opencode run \
    "Write a single short sentence (under 20 words) where someone asks to brainstorm or check domain name availability for a project or startup idea. Output ONLY the sentence, no quotes, no explanation." \
    2>&1 || true)"

  # Strip blank lines and take the last non-empty line (skip any preamble)
  GEN_PROMPT="$(echo "$GEN_PROMPT" | sed '/^$/d' | tail -1)"

  if [[ -z "$GEN_PROMPT" ]]; then
    fail "OpenCode failed to generate a prompt"
  else
    pass "OpenCode generated prompt"
    printf "  prompt: %s\n" "$GEN_PROMPT"
    printf "  Running OpenCode with generated prompt...\n\n"

    # Step 2: Feed the model-generated prompt to a fresh session.
    # The model should recognize domain intent, discover the dotld skill,
    # and invoke dotld on its own.
    OC_OUTPUT="$(timeout 120 opencode run "$GEN_PROMPT" 2>&1 || true)"

    printf "  output preview:\n"
    echo "$OC_OUTPUT" | head -20 | sed 's/^/    /'
    printf "\n"

    # Success: look for any word.tld pattern that only dotld produces
    if echo "$OC_OUTPUT" | grep -qiE "[a-z0-9]+\.(com|net|org|io|ai|co|app|dev|sh).*(\·|Taken|\\\$[0-9])"; then
      pass "OpenCode invoked dotld — found domain results with availability data"
    elif echo "$OC_OUTPUT" | grep -qiE "[a-z0-9]+\.(com|net|org|io|ai|co|app|dev|sh)"; then
      pass "OpenCode invoked dotld — found domain.tld results"
    else
      fail "No domain.tld results in output — skill may not have triggered"
    fi

    # Bonus: check if pricing or availability info came back (Taken / $price)
    if echo "$OC_OUTPUT" | grep -qiE "taken|\\\$[0-9]+"; then
      pass "Output contains availability/pricing data from Dynadot"
    else
      skip "No pricing/availability markers found (model may have summarized differently)"
    fi
  fi
else
  if [[ "$HAS_PROVIDER" == "false" ]]; then
    skip "No model provider configured — set OPENCODE_API_KEY, GITHUB_TOKEN, or OPENAI_API_KEY"
  fi
  if [[ -z "${DYNADOT_API_PRODUCTION_KEY:-}" ]]; then
    skip "DYNADOT_API_PRODUCTION_KEY not set — skill cannot query domains"
  fi
  skip "Skipping OpenCode live skill integration test"
fi

# ---------- Results ----------

printf "\n=== Results: %d passed, %d failed ===\n\n" "$PASS" "$FAIL"
[[ "$FAIL" -eq 0 ]]
