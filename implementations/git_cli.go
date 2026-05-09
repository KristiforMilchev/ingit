package implementations

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"ingit/models"
)

type GitCLI struct{}

func NewGitCLI() *GitCLI {
	return &GitCLI{}
}

func (g *GitCLI) DiscoverRepos() []models.Repo {
	cwd, err := os.Getwd()
	if err != nil {
		return nil
	}

	cwdRoot := g.gitRoot(cwd)

	var childRepos []models.Repo
	seen := map[string]bool{}

	entries, err := os.ReadDir(cwd)
	if err == nil {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			path := filepath.Join(cwd, entry.Name())
			root := g.gitRoot(path)
			if root == "" {
				continue
			}

			if root == cwdRoot {
				continue
			}

			if seen[root] {
				continue
			}

			seen[root] = true
			childRepos = append(childRepos, models.Repo{
				Name: filepath.Base(root),
				Path: root,
			})
		}
	}

	if len(childRepos) > 0 {
		return childRepos
	}

	if cwdRoot != "" {
		return []models.Repo{{
			Name: filepath.Base(cwdRoot),
			Path: cwdRoot,
		}}
	}

	return nil
}

func (g *GitCLI) gitRoot(path string) string {
	out, err := g.git(path, "rev-parse", "--show-toplevel")
	if err != nil {
		return ""
	}

	return strings.TrimSpace(out)
}

func (g *GitCLI) ReadGitState(repo models.Repo) models.GitState {
	out, err := g.git(repo.Path, "status", "--porcelain=v1", "-b")
	if err != nil {
		return models.GitState{Error: err.Error()}
	}

	state := models.GitState{}
	state.Graph = g.ReadBranchGraph(repo)

	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "## ") {
			parseBranchLine(&state, strings.TrimPrefix(line, "## "))
			continue
		}
		if len(line) < 4 {
			continue
		}

		index := string(line[0])
		worktree := string(line[1])
		path := strings.TrimSpace(line[3:])
		if strings.Contains(path, " -> ") {
			parts := strings.Split(path, " -> ")
			path = parts[len(parts)-1]
		}

		state.Files = append(state.Files, models.FileStatus{
			Path:     path,
			Index:    index,
			Worktree: worktree,
		})
	}

	return state
}

func (g *GitCLI) ReadBranchGraph(repo models.Repo) string {
	out, err := g.git(repo.Path,
		"log",
		"--graph",
		"--oneline",
		"--decorate",
		"--all",
		"--date-order",
		"-n",
		"60",
	)
	if err != nil {
		return err.Error()
	}
	if strings.TrimSpace(out) == "" {
		return "No branch graph available."
	}
	return out
}

func (g *GitCLI) ReadRecentCommits(repo models.Repo) []models.GitCommit {
	out, err := g.git(repo.Path,
		"log",
		"--decorate",
		"--oneline",
		"-n",
		"8",
	)
	if err != nil {
		return nil
	}

	var commits []models.GitCommit
	for _, line := range strings.Split(strings.TrimSpace(out), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		commit := models.GitCommit{Hash: parts[0]}
		if len(parts) > 1 {
			rest := parts[1]
			if strings.HasPrefix(rest, "(") {
				end := strings.Index(rest, ")")
				if end != -1 {
					commit.Refs = rest[:end+1]
					commit.Message = strings.TrimSpace(rest[end+1:])
				} else {
					commit.Message = rest
				}
			} else {
				commit.Message = rest
			}
		}
		commits = append(commits, commit)
	}

	return commits
}

func parseBranchLine(state *models.GitState, line string) {
	parts := strings.Split(line, "...")
	state.Branch = parts[0]
	if strings.Contains(line, "[ahead ") {
		state.Ahead = extractCount(line, "ahead ")
	}
	if strings.Contains(line, "behind ") {
		state.Behind = extractCount(line, "behind ")
	}
}

func extractCount(s, marker string) int {
	idx := strings.Index(s, marker)
	if idx == -1 {
		return 0
	}
	start := idx + len(marker)
	end := start
	for end < len(s) && s[end] >= '0' && s[end] <= '9' {
		end++
	}
	var n int
	_, _ = fmt.Sscanf(s[start:end], "%d", &n)
	return n
}

func (g *GitCLI) ReadDiffForFile(repo models.Repo, file models.FileStatus) string {
	unstaged := g.readDiff(repo, file.Path, false)
	if isRealDiff(unstaged) {
		return unstaged
	}

	staged := g.readDiff(repo, file.Path, true)
	if isRealDiff(staged) {
		return staged
	}

	return unstaged
}

func isRealDiff(value string) bool {
	return strings.HasPrefix(value, "diff --git") || strings.Contains(value, "\n@@")
}

func (g *GitCLI) readDiff(repo models.Repo, path string, staged bool) string {
	args := []string{"diff", "--no-color", "--", path}
	if staged {
		args = []string{"diff", "--cached", "--no-color", "--", path}
	}

	out, err := g.git(repo.Path, args...)
	if err != nil {
		return err.Error()
	}

	if strings.TrimSpace(out) == "" {
		if staged {
			return "No staged diff for selected file."
		}
		return "No unstaged diff for selected file."
	}
	return out
}

func (g *GitCLI) StageFile(repo models.Repo, path string) error {
	_, err := g.git(repo.Path, "add", "--", path)
	return err
}

func (g *GitCLI) UnstageFile(repo models.Repo, path string) error {
	_, err := g.git(repo.Path, "restore", "--staged", "--", path)
	return err
}

func (g *GitCLI) StageAll(repo models.Repo) error {
	_, err := g.git(repo.Path, "add", "-A")
	return err
}

func (g *GitCLI) UnstageAll(repo models.Repo) error {
	_, err := g.git(repo.Path, "restore", "--staged", ".")
	return err
}

func (g *GitCLI) ExecutePushPlan(repo models.Repo, plan models.PushPlanExecution) error {
	if strings.TrimSpace(plan.Branch) == "" {
		return fmt.Errorf("branch name is required")
	}
	if strings.TrimSpace(plan.Message) == "" {
		return fmt.Errorf("commit message is required for %s", plan.Branch)
	}
	if len(plan.Files) == 0 {
		return fmt.Errorf("no files assigned to %s", plan.Branch)
	}

	base := strings.TrimSpace(plan.Base)
	if base == "" || base == "current" {
		base = "HEAD"
	}

	currentBranchOut, _ := g.git(repo.Path, "rev-parse", "--abbrev-ref", "HEAD")
	currentBranch := strings.TrimSpace(currentBranchOut)

	_, err := g.git(repo.Path, "checkout", base)
	if err != nil {
		return fmt.Errorf("checkout base %s: %w", base, err)
	}

	_, err = g.git(repo.Path, "checkout", "-B", plan.Branch)
	if err != nil {
		_ = g.restoreBranch(repo, currentBranch)
		return fmt.Errorf("create branch %s: %w", plan.Branch, err)
	}

	_, _ = g.git(repo.Path, "restore", "--staged", ".")

	args := append([]string{"add", "--"}, plan.Files...)
	_, err = g.git(repo.Path, args...)
	if err != nil {
		_ = g.restoreBranch(repo, currentBranch)
		return fmt.Errorf("add files for %s: %w", plan.Branch, err)
	}

	_, err = g.git(repo.Path, "commit", "-m", plan.Message)
	if err != nil {
		_ = g.restoreBranch(repo, currentBranch)
		return fmt.Errorf("commit %s: %w", plan.Branch, err)
	}

	_, err = g.git(repo.Path, "push", "-u", "origin", plan.Branch)
	if err != nil {
		_ = g.restoreBranch(repo, currentBranch)
		return fmt.Errorf("push %s: %w", plan.Branch, err)
	}

	_ = g.restoreBranch(repo, currentBranch)
	return nil
}

func (g *GitCLI) ExecuteQuickPush(repo models.Repo, message string) error {
	message = strings.TrimSpace(message)
	if message == "" {
		return fmt.Errorf("commit message is required")
	}

	branchOut, err := g.git(repo.Path, "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return fmt.Errorf("read current branch: %w", err)
	}

	branch := strings.TrimSpace(branchOut)
	if branch == "" || branch == "HEAD" {
		return fmt.Errorf("cannot quick push from detached HEAD")
	}

	_, err = g.git(repo.Path, "add", "-A")
	if err != nil {
		return fmt.Errorf("stage all changes: %w", err)
	}

	_, err = g.git(repo.Path, "commit", "-m", message)
	if err != nil {
		return fmt.Errorf("commit current branch: %w", err)
	}

	_, err = g.git(repo.Path, "push")
	if err == nil {
		return nil
	}

	_, err = g.git(repo.Path, "push", "-u", "origin", branch)
	if err != nil {
		return fmt.Errorf("push current branch %s: %w", branch, err)
	}

	return nil
}

func (g *GitCLI) restoreBranch(repo models.Repo, branch string) error {
	if strings.TrimSpace(branch) == "" || branch == "HEAD" {
		return nil
	}
	_, err := g.git(repo.Path, "checkout", branch)
	return err
}

func (g *GitCLI) git(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return stdout.String(), fmt.Errorf("%s", msg)
	}
	return stdout.String(), nil
}
