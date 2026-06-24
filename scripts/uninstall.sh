#!/usr/bin/env bash
# Remove mta on macOS/Linux: stops + removes the background service and deletes
# the binary (delegates to `mta self uninstall`, which knows the per-OS teardown).
#
#   curl -fsSL https://raw.githubusercontent.com/simtabi/ms-teams-activity/main/scripts/uninstall.sh | sh
#   ./scripts/uninstall.sh            # remove service + binary
#   ./scripts/uninstall.sh --purge    # also delete config and runtime data
set -eu

PURGE=0
for a in "$@"; do
  case "$a" in
    --purge) PURGE=1 ;;
    *) echo "unknown option: $a" >&2; exit 1 ;;
  esac
done

mta="$(command -v mta || true)"
if [ -z "$mta" ]; then
  for p in "$HOME/.local/bin/mta" /usr/local/bin/mta; do
    [ -x "$p" ] && mta="$p" && break
  done
fi
if [ -z "$mta" ]; then
  echo "mta not found on PATH or in the usual install dirs; nothing to do." >&2
  exit 1
fi

cmd="\"$mta\" self uninstall --yes"
[ "$PURGE" = "1" ] && cmd="$cmd --purge"
echo "Removing mta (service + binary)..."
eval "$cmd"
echo "Done."
