#!/usr/bin/env bash
#
# post-gen.sh — runs after OpenAPI code generation for generic Python SDKs.
# Validates that docs/CHANGELOG.md was not lost during generation.
#

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../../.." && pwd)"

CHANGELOG_PATH="${REPO_ROOT}/docs/CHANGELOG.md"

# ---------------------------------------------------------------------------
# Guard: ensure docs/CHANGELOG.md exists after generation
# ---------------------------------------------------------------------------
if [[ ! -f "${CHANGELOG_PATH}" ]]; then
    if [[ "${ALLOW_MISSING_CHANGELOG:-false}" == "true" ]]; then
        echo "[post-gen] WARNING: docs/CHANGELOG.md is missing after generation (allowed by ALLOW_MISSING_CHANGELOG=true)"
        echo "[post-gen] This is expected for first-ever releases. A CHANGELOG.md should be created manually."
    else
        echo "[post-gen] ERROR: docs/CHANGELOG.md is missing after generation!"
        echo "[post-gen] The changelog may have been deleted during the generation process."
        echo "[post-gen] If this is a first release with no prior changelog, set ALLOW_MISSING_CHANGELOG=true"
        exit 1
    fi
else
    echo "[post-gen] docs/CHANGELOG.md exists — OK"
fi
