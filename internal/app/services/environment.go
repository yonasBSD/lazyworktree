package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BuildCommandEnv builds environment variables for worktree commands.
func BuildCommandEnv(branch, wtPath, repoKey, mainWorktreePath string) map[string]string {
	owner, repo := SplitRepoKey(repoKey)
	return map[string]string{
		"WORKTREE_BRANCH":    branch,
		"MAIN_WORKTREE_PATH": mainWorktreePath,
		"WORKTREE_PATH":      wtPath,
		"WORKTREE_NAME":      filepath.Base(wtPath),
		"REPO_NAME":          repoKey,
		"REPO_OWNER":         owner,
		"REPO_REPONAME":      repo,
	}
}

// SplitRepoKey splits a repository key like "owner/repo" into owner and repo.
// If the key contains no slash, owner is empty and repo is the full key.
func SplitRepoKey(repoKey string) (owner, repo string) {
	if i := strings.IndexByte(repoKey, '/'); i >= 0 {
		return repoKey[:i], repoKey[i+1:]
	}
	return "", repoKey
}

// ExpandWithEnv expands environment variables using the provided map first.
func ExpandWithEnv(input string, env map[string]string) string {
	if input == "" {
		return ""
	}
	return os.Expand(input, func(key string) string {
		if val, ok := env[key]; ok {
			return val
		}
		return os.Getenv(key)
	})
}

// EnvMapToList converts environment variables to KEY=VALUE pairs.
func EnvMapToList(env map[string]string) []string {
	if len(env) == 0 {
		return nil
	}
	out := make([]string, 0, len(env))
	for key, val := range env {
		out = append(out, fmt.Sprintf("%s=%s", key, val))
	}
	return out
}
