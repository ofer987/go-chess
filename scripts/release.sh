#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 2 ]]; then
    echo "Usage: $0 <version> <PR#> [PR#...]" >&2
    echo "Example: $0 1.0.0 42 57 61" >&2
    exit 1
fi

VERSION=$1
shift
PRS=("$@")

# Collect merge SHAs with their merged-at timestamps, then sort ascending by date.
declare -a ENTRIES=()
for PR in "${PRS[@]}"; do
    DATA=$(gh pr view "$PR" --json number,mergeCommit,mergedAt \
        --jq '[.number, .mergeCommit.oid, .mergedAt] | @tsv')
    SHA=$(echo "$DATA" | cut -f2)
    MERGED_AT=$(echo "$DATA" | cut -f3)

    if [[ -z "$SHA" || "$SHA" == "null" ]]; then
        echo "PR #$PR has no merge commit — was it merged? Skipping." >&2
        continue
    fi

    ENTRIES+=("$MERGED_AT $SHA PR#$PR")
done

if [[ ${#ENTRIES[@]} -eq 0 ]]; then
    echo "No merged PRs found — nothing to do." >&2
    exit 1
fi

# Sort ascending by merged-at timestamp (ISO-8601 sorts lexicographically).
mapfile -t SORTED < <(printf '%s\n' "${ENTRIES[@]}" | sort)

git checkout master
git pull origin master
git checkout -b "releases/$VERSION"

for ENTRY in "${SORTED[@]}"; do
    SHA=$(echo "$ENTRY" | awk '{print $2}')
    LABEL=$(echo "$ENTRY" | awk '{print $3}')
    DATE=$(echo "$ENTRY" | awk '{print $1}')
    echo "Cherry-picking $LABEL ($SHA, merged $DATE)"
    git cherry-pick -m 1 "$SHA"
done

git push origin "releases/$VERSION"
echo "Release branch releases/$VERSION pushed."
