#!/usr/bin/env bash
# Read-only host survey. It intentionally changes no package, service, firewall,
# container, route, SELinux setting or file.
set -u

section() { printf '\n===== %s =====\n' "$1"; }
run() {
  printf '\n$ %s\n' "$*"
  "$@" 2>&1 || printf '[WARN] command failed with code %s\n' "$?"
}

section "Identity"
run date --iso-8601=seconds
run hostnamectl
run uname -a
run cat /etc/os-release
run ip -br address

section "Capacity"
run uptime
run free -h
run df -hT /
run df -ih /

section "Existing failed units"
run systemctl --failed --no-pager

section "Important host services (snapshot only)"
for unit in docker.service containerd.service firewalld.service fail2ban.service netdata.service sshd.service; do
  printf '%-28s ' "$unit"
  systemctl is-active "$unit" 2>/dev/null || true
  systemctl is-enabled "$unit" 2>/dev/null | sed 's/^/  enabled: /' || true
done

section "Docker snapshot"
if command -v docker >/dev/null 2>&1; then
  run docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}'
  run docker network ls
else
  echo "docker command not found"
fi

section "Listening sockets"
run ss -lntup

section "Firewall snapshot"
if command -v firewall-cmd >/dev/null 2>&1; then
  run firewall-cmd --state
  run firewall-cmd --get-active-zones
  run firewall-cmd --list-all
else
  echo "firewall-cmd not found"
fi

section "SELinux"
run getenforce
run sestatus

section "Potential somonwatch conflicts"
for path in /opt/somonwatch /etc/somonwatch /var/lib/somonwatch /etc/systemd/system/somonwatch.service; do
  if [[ -e "$path" ]]; then
    echo "[NOTICE] already exists: $path"
    ls -ld "$path"
  else
    echo "[OK] absent: $path"
  fi
done
if systemctl list-unit-files --no-legend 2>/dev/null | grep -q '^somonwatch.service'; then
  run systemctl status somonwatch.service --no-pager
fi

section "Build prerequisites"
for cmd in go gcc sqlite3; do
  if command -v "$cmd" >/dev/null 2>&1; then
    echo "[OK] $cmd: $(command -v "$cmd")"
    "$cmd" --version 2>&1 | head -n 2 || true
  else
    echo "[MISSING] $cmd"
  fi
done
if [[ -f /usr/include/sqlite3.h ]]; then
  echo "[OK] /usr/include/sqlite3.h"
else
  echo "[MISSING] sqlite3.h (package sqlite-devel)"
fi

section "Outbound HTTPS checks"
if command -v curl >/dev/null 2>&1; then
  curl -sS -L --max-time 20 -o /dev/null -w 'Somon: HTTP %{http_code}, %{size_download} bytes, %{time_total}s\n' \
    'https://somon.tj/nedvizhimost/arenda-kvartir/dushanbe/' || echo '[WARN] Somon HTTPS check failed'
  curl -sS -L --max-time 20 -o /dev/null -w 'Telegram API: HTTP %{http_code}, %{time_total}s\n' \
    'https://api.telegram.org/' || echo '[WARN] Telegram HTTPS check failed'
else
  echo "curl not found"
fi

section "Conclusion"
echo "This script made no changes. Save this output and compare Docker/services/sockets after deployment."
echo "somonwatch requires no inbound port and must not require Traefik or firewalld changes."
