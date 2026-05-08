Create a patch release branch by cherry-picking the merge commits of the specified GitHub Pull Requests onto a new `releases/x.y.z` branch.

**Arguments:** `$ARGUMENTS` — a semver version (e.g. `1.2.3`) followed by one or more PR numbers (e.g. `42 57 61`), or a natural language description like "patch release 1.2.3 with PRs #42, #57, and #61".

## Steps

1. Parse the version and PR numbers from the arguments.
2. For each PR number, fetch its merge commit SHA and merged-at timestamp:

   ```bash
   gh pr view <PR#> --json number,title,mergeCommit,mergedAt,state
   ```

   - If the PR is not merged (no merge commit), warn the user and skip it.
   - Collect: PR number, title, merge commit SHA, mergedAt timestamp.

3. Sort the collected PRs by `mergedAt` ascending (oldest merge first), so history lands in chronological order on the release branch.
4. Show the user the sorted list of PRs that will be cherry-picked and ask for confirmation before proceeding.
5. Create the release branch from master:

   ```bash
   git checkout master
   git pull origin master
   git checkout -b releases/<version>
   ```

6. Cherry-pick each merge commit in sorted order:

   ```bash
   git cherry-pick -m 1 <merge-sha>
   ```

   - If a conflict occurs, stop and tell the user exactly which PR caused the conflict and what files are conflicting. Do not attempt to resolve conflicts automatically.

7. Run tests to verify the release branch is healthy:

   ```bash
   go test ./...
   ```

8. Push the release branch:

   ```bash
   git push origin releases/<version>
   ```

9. Report a summary: version, branch name, and the list of PRs (number + title) that were included.
