#!/usr/bin/env bash
set -Eeuo pipefail

if (( EUID != 0 )); then
  echo "ERROR: run as root" >&2
  exit 1
fi
command -v sqlite3 >/dev/null || { echo "ERROR: sqlite3 CLI is required" >&2; exit 1; }

STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
DEST="${1:-/var/backups/somonwatch/manual-${STAMP}}"
install -d -o root -g root -m 0700 "$DEST"

for path in /opt/somonwatch/somonwatch /opt/somonwatch/README.md /opt/somonwatch/RUNBOOK_ALMALINUX_9.md \
            /etc/somonwatch/somonwatch.env /etc/systemd/system/somonwatch.service; do
  [[ -e "$path" ]] && cp -a "$path" "$DEST/"
done
if [[ -f /var/lib/somonwatch/somonwatch.db ]]; then
  sqlite3 /var/lib/somonwatch/somonwatch.db ".timeout 5000" ".backup '$DEST/somonwatch.db'"
  chmod 0600 "$DEST/somonwatch.db"
fi
sha256sum "$DEST"/* > "$DEST/SHA256SUMS" 2>/dev/null || true
chmod 0600 "$DEST/SHA256SUMS" 2>/dev/null || true

echo "$DEST"
