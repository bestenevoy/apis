#!/usr/bin/env bash
set -euo pipefail

# Roll back to the previous release on the server.
# Usage: bash scripts/rollback.sh [tag]
# - no tag: switches to the second newest release in /opt/wrzapi/releases
# - tag:    switches to /opt/wrzapi/releases/wrzapi_<tag>

RELEASE_DIR=/opt/wrzapi/releases
CURRENT_LINK=/opt/wrzapi/current

if [[ ${1-} != "" ]]; then
  TARGET="$RELEASE_DIR/wrzapi_$1"
  if [[ ! -f "$TARGET" ]]; then
    echo "release not found: $TARGET" >&2
    exit 1
  fi
else
  TARGET=$(ls -1t "$RELEASE_DIR"/wrzapi_* 2>/dev/null | sed -n '2p')
  if [[ -z "${TARGET}" ]]; then
    echo "no previous release found" >&2
    exit 1
  fi
fi

ln -sfn "$TARGET" "$CURRENT_LINK"
sudo systemctl restart wrzapi
sudo systemctl status wrzapi --no-pager