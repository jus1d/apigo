#!/usr/bin/env bash
set -euo pipefail

OLD="apigo"

if [ $# -ne 1 ]; then
    echo "Usage: $0 <new-name>"
    exit 1
fi

NEW="$1"

if [ "$OLD" = "$NEW" ]; then
    echo "New name is the same as current name, nothing to do."
    exit 0
fi

# Replace OLD in all tracked files (excludes .git, vendor via .gitignore, etc.)
git ls-files -z | xargs -0 grep -lZ -F "${OLD}" 2>/dev/null | xargs -0 sed -i '' "s|${OLD}|${NEW}|g"

echo "Renamed module from '${OLD}' to '${NEW}'"
