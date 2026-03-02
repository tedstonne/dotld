#!/usr/bin/env bash
set -euo pipefail

SKILL_DIR="$HOME/.claude/skills/dotld"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

if [[ "${1:-}" == "clean" ]]; then
  rm -rf "$SKILL_DIR"
  printf "Removed %s\n" "$SKILL_DIR"
  exit 0
fi

mkdir -p "$SKILL_DIR"
cp -r "$PROJECT_DIR/skills/dotld/"* "$SKILL_DIR/"

printf "Skill installed to %s\n\n" "$SKILL_DIR"
printf "To test:\n"
printf "  1. Start a new Claude Code session\n"
printf "  2. Type: /dotld acme\n"
printf "  3. Verify it runs dotld and shows domain results\n\n"
printf "To clean up:\n"
printf "  bash %s clean\n" "$0"
