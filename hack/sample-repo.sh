#!/usr/bin/env bash
set -euo pipefail

REPO_DIR="/tmp/repo"
WORKTREE_BASE="/tmp/rworktrees-throwaway"

# ── Cleanup ────────────────────────────────────────────────────────────────────
echo "🧹 Cleaning up existing repo/worktrees..."
rm -rf "$REPO_DIR" "$WORKTREE_BASE"

# ── Init ───────────────────────────────────────────────────────────────────────
echo "📁 Initialising repo at $REPO_DIR..."
mkdir -p "$REPO_DIR"
git -C "$REPO_DIR" init -q
git -C "$REPO_DIR" config user.email "throwaway@local"
git -C "$REPO_DIR" config user.name "Throwaway"

# ── Helper ─────────────────────────────────────────────────────────────────────
commit() { # commit <message> [file]
  local msg="$1"
  local file="${2:-file_$RANDOM.txt}"
  echo "$msg — $(date +%s%N)" >"$REPO_DIR/$file"
  git -C "$REPO_DIR" add -A
  git -C "$REPO_DIR" commit -q -m "$msg"
}

# ── main branch ───────────────────────────────────────────────────────────────
echo "✏️  Creating commits on main..."
git -C "$REPO_DIR" checkout -q -b main 2>/dev/null || true

commit "chore: initial scaffold" "README.md"
commit "feat: add config file" "config.yaml"
commit "feat: add main module" "main.py"
commit "fix: correct typo in README" "README.md"
commit "docs: expand config comments" "config.yaml"
commit "refactor: restructure main" "main.py"
commit "chore: add .gitignore" ".gitignore"
commit "test: add basic test suite" "test_main.py"
commit "ci: add github actions stub" ".github_actions.yml"
commit "release: v1.0.0" "CHANGELOG.md"

# ── feature branches ──────────────────────────────────────────────────────────
echo "🌿 Creating feature branches..."

git -C "$REPO_DIR" checkout -q -b feature/auth
commit "feat(auth): scaffold auth module" "auth.py"
commit "feat(auth): add JWT support" "auth.py"
commit "test(auth): add auth unit tests" "test_auth.py"
commit "fix(auth): handle token expiry" "auth.py"

git -C "$REPO_DIR" checkout -q main
git -C "$REPO_DIR" checkout -q -b feature/api
commit "feat(api): add REST endpoints" "api.py"
commit "feat(api): add pagination" "api.py"
commit "test(api): endpoint smoke tests" "test_api.py"

git -C "$REPO_DIR" checkout -q main
git -C "$REPO_DIR" checkout -q -b feature/db
commit "feat(db): add ORM models" "models.py"
commit "feat(db): add migrations" "migrations.sql"
commit "fix(db): fix FK constraint" "models.py"

git -C "$REPO_DIR" checkout -q main
git -C "$REPO_DIR" checkout -q -b bugfix/memory-leak
commit "fix: patch memory leak in loop" "main.py"
commit "fix: release handles on exit" "main.py"

git -C "$REPO_DIR" checkout -q main
git -C "$REPO_DIR" checkout -q -b release/v2.0
commit "chore: bump version to 2.0" "version.txt"
commit "docs: update CHANGELOG for v2" "CHANGELOG.md"

git -C "$REPO_DIR" checkout -q main
git -C "$REPO_DIR" checkout -q -b experiment/new-engine
commit "wip: prototype new engine" "engine.py"
commit "wip: rough benchmark harness" "bench.py"

git -C "$REPO_DIR" checkout -q main

# ── worktrees ─────────────────────────────────────────────────────────────
echo "🗂️  Creating worktrees..."
mkdir -p "$WORKTREE_BASE"

git -C "$REPO_DIR" worktree add -q "$WORKTREE_BASE/wt-auth" feature/auth
git -C "$REPO_DIR" worktree add -q "$WORKTREE_BASE/wt-api" feature/api
git -C "$REPO_DIR" worktree add -q "$WORKTREE_BASE/wt-db" feature/db
git -C "$REPO_DIR" worktree add -q "$WORKTREE_BASE/wt-bugfix" bugfix/memory-leak
git -C "$REPO_DIR" worktree add -q "$WORKTREE_BASE/wt-release" release/v2.0

# ── Untracked and staged files in main repo ────────────────────────────────
echo "📝 Adding random files to main repo..."
echo "Untracked change 1" >"$REPO_DIR/untracked_main_1.txt"
echo "Untracked change 2" >"$REPO_DIR/untracked_main_2.txt"
echo "Staged change" >"$REPO_DIR/staged_main.txt"
git -C "$REPO_DIR" add staged_main.txt

# ── Untracked and staged files in worktrees ────────────────────────────────
echo "📝 Adding random files to worktrees..."

# wt-auth
echo "Untracked in auth worktree" >"$WORKTREE_BASE/wt-auth/untracked_auth.txt"
echo "Staged in auth worktree" >"$WORKTREE_BASE/wt-auth/staged_auth.txt"
git -C "$WORKTREE_BASE/wt-auth" add staged_auth.txt

# wt-api
echo "Untracked in api worktree" >"$WORKTREE_BASE/wt-api/untracked_api.txt"
echo "Staged in api worktree" >"$WORKTREE_BASE/wt-api/staged_api.txt"
git -C "$WORKTREE_BASE/wt-api" add staged_api.txt

# wt-db
echo "Untracked in db worktree" >"$WORKTREE_BASE/wt-db/untracked_db.txt"
echo "Staged in db worktree" >"$WORKTREE_BASE/wt-db/staged_db.txt"
git -C "$WORKTREE_BASE/wt-db" add staged_db.txt

# ── Summary ───────────────────────────────────────────────────────────────
echo ""
echo "✅ Done!"
echo ""
echo "  Repo       : $REPO_DIR"
echo "  Worktrees  : $WORKTREE_BASE/"
echo ""
echo "  Branches:"
git -C "$REPO_DIR" branch | sed 's/^/    /'
echo ""
echo "  Worktrees:"
git -C "$REPO_DIR" worktree list | sed 's/^/    /'
echo ""
echo "  Commits (main):"
git -C "$REPO_DIR" log --oneline main | sed 's/^/    /'
