#!/usr/bin/env bash
set -Eeuo pipefail

ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)"
START=0
if [[ "${1:-}" == "--start" ]]; then
  START=1
elif [[ $# -gt 0 ]]; then
  echo "Usage: sudo $0 [--start]" >&2
  exit 2
fi

if (( EUID != 0 )); then
  echo "ERROR: run as root" >&2
  exit 1
fi
if [[ ! -x "$ROOT_DIR/dist/somonwatch" ]]; then
  echo "ERROR: dist/somonwatch not found; run scripts/build.sh first" >&2
  exit 1
fi
if systemctl is-active --quiet somonwatch.service 2>/dev/null; then
  echo "ERROR: somonwatch.service is already active. Stop only this unit before an upgrade:" >&2
  echo "  systemctl stop somonwatch.service" >&2
  exit 1
fi

if ! getent group somonwatch >/dev/null; then
  groupadd --system somonwatch
fi
if ! id somonwatch >/dev/null 2>&1; then
  useradd --system --gid somonwatch --home-dir /var/lib/somonwatch --shell /sbin/nologin --comment "Somon Rent Watcher" somonwatch
fi

install -d -o root -g root -m 0755 /opt/somonwatch
install -d -o root -g root -m 0700 /etc/somonwatch
install -d -o somonwatch -g somonwatch -m 0700 /var/lib/somonwatch /var/lib/somonwatch/debug

STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
BACKUP_DIR="/var/backups/somonwatch/install-${STAMP}"
mkdir -p "$BACKUP_DIR"
for path in /opt/somonwatch/somonwatch /etc/systemd/system/somonwatch.service /etc/somonwatch/somonwatch.env; do
  if [[ -e "$path" ]]; then
    cp -a "$path" "$BACKUP_DIR/"
  fi
done
chmod 0700 "$BACKUP_DIR"

install -o root -g root -m 0755 "$ROOT_DIR/dist/somonwatch" /opt/somonwatch/somonwatch
install -o root -g root -m 0644 "$ROOT_DIR/deploy/somonwatch.service" /etc/systemd/system/somonwatch.service
install -o root -g root -m 0644 "$ROOT_DIR/docs/RUNBOOK_ALMALINUX_9.md" /opt/somonwatch/RUNBOOK_ALMALINUX_9.md
install -o root -g root -m 0644 "$ROOT_DIR/README.md" /opt/somonwatch/README.md

if [[ ! -e /etc/somonwatch/somonwatch.env ]]; then
  install -o root -g root -m 0600 "$ROOT_DIR/deploy/somonwatch.env.example" /etc/somonwatch/somonwatch.env
  echo "Created /etc/somonwatch/somonwatch.env with placeholders."
else
  chmod 0600 /etc/somonwatch/somonwatch.env
  chown root:root /etc/somonwatch/somonwatch.env
  echo "Preserved existing /etc/somonwatch/somonwatch.env."
fi

chown -R somonwatch:somonwatch /var/lib/somonwatch
chmod 0700 /var/lib/somonwatch /var/lib/somonwatch/debug
if [[ -e /var/lib/somonwatch/somonwatch.db ]]; then
  chmod 0600 /var/lib/somonwatch/somonwatch.db*
fi

if command -v restorecon >/dev/null 2>&1; then
  restorecon -RF /opt/somonwatch /etc/somonwatch /var/lib/somonwatch /etc/systemd/system/somonwatch.service || true
fi
systemctl daemon-reload

if (( START == 1 )); then
  if grep -q 'replace_me' /etc/somonwatch/somonwatch.env; then
    echo "ERROR: environment still contains placeholders; service was not started" >&2
    exit 1
  fi
  systemctl enable --now somonwatch.service
  systemctl --no-pager --full status somonwatch.service || true
else
  echo "Installed but deliberately not enabled or started."
  echo "Next: edit /etc/somonwatch/somonwatch.env, run doctor, then:"
  echo "  systemctl enable --now somonwatch.service"
fi

echo "Backup of replaced install files: $BACKUP_DIR"
echo "No Docker container, Traefik rule, firewall port, route or SELinux boolean was changed."
