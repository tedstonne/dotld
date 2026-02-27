#!/usr/bin/env bash
set -euo pipefail

ITEM="dotld"
ACCOUNT="sasso.1password.com"

DYNADOT_API_PRODUCTION_KEY=$(op read "op://Shared/$ITEM/DYNADOT_API_PRODUCTION_KEY" --account "$ACCOUNT")

cat > .env <<EOF
DYNADOT_API_PRODUCTION_KEY=$DYNADOT_API_PRODUCTION_KEY
EOF

echo "Generated .env from 1Password ($ITEM)"
