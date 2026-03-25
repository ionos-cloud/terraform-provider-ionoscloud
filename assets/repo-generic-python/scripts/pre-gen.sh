#!/usr/bin/env bash
#
# pre-gen.sh — runs before OpenAPI code generation for generic Python SDKs.
# Cleans the docs/ directory while preserving docs/CHANGELOG.md.
#

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../../.." && pwd)"

CHANGELOG_PATH="${REPO_ROOT}/docs/CHANGELOG.md"
CHANGELOG_BACKUP="/tmp/_changelog_backup_$$"

# ---------------------------------------------------------------------------
# 1. Backup docs/CHANGELOG.md before wiping docs/
# ---------------------------------------------------------------------------
if [[ -f "${CHANGELOG_PATH}" ]]; then
    echo "[pre-gen] Backing up docs/CHANGELOG.md"
    cp "${CHANGELOG_PATH}" "${CHANGELOG_BACKUP}"
else
    echo "[pre-gen] WARNING: docs/CHANGELOG.md not found — nothing to back up"
fi

# ---------------------------------------------------------------------------
# 2. Remove generated docs (the original destructive step)
# ---------------------------------------------------------------------------
echo "[pre-gen] Removing docs/ directory"
rm -rf "${REPO_ROOT}/docs/"

# ---------------------------------------------------------------------------
# 3. Restore docs/CHANGELOG.md from backup
# ---------------------------------------------------------------------------
if [[ -f "${CHANGELOG_BACKUP}" ]]; then
    mkdir -p "${REPO_ROOT}/docs"
    cp "${CHANGELOG_BACKUP}" "${CHANGELOG_PATH}"
    rm -f "${CHANGELOG_BACKUP}"
    echo "[pre-gen] Restored docs/CHANGELOG.md"
else
    echo "[pre-gen] No changelog backup to restore (first release?)"
fi
