package controllers

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"ingit/interfaces"
	"ingit/models"
	"ingit/states"
)

type BranchesModel struct {
	State states.BranchesState
	Git   interfaces.GitProvider
}

func NewBranchesModel(git interfaces.GitProvider) BranchesModel {
	return BranchesModel{
		State: states.BranchesState{},
		Git:   git,
	}
}

func (m *BranchesModel) Load(repoPath string) {
	m.State.RepoPath = repoPath
	m.State.Error = ""
	m.State.Done = false
	m.State.CheckedOut = false
	m.State.Branches = nil
	m.State.Commits = nil
	m.State.Files = nil
	m.State.Selected = 0

	current, _ := m.Git.GitOutput(repoPath, "branch", "--show-current")
	m.State.CurrentBranch = strings.TrimSpace(current)

	out, err := m.Git.GitOutput(repoPath, "branch", "--all", "--no-color")
	if err != nil {
		m.State.Error = cleanGitError(out, err)
		return
	}

	for _, line := range strings.Split(out, "\n") {
		raw := strings.TrimSpace(line)
		if raw == "" {
			continue
		}

		current := strings.HasPrefix(raw, "*")
		raw = strings.TrimSpace(strings.TrimPrefix(raw, "*"))

		if strings.Contains(raw, "HEAD ->") {
			continue
		}

		remote := strings.HasPrefix(raw, "remotes/")
		name := strings.TrimPrefix(raw, "remotes/origin/")

		if name == "" || name == "HEAD" {
			continue
		}

		if containsBranch(m.State.Branches, name) {
			continue
		}

		m.State.Branches = append(m.State.Branches, models.GitBranch{
			Name:    name,
			Current: current || name == m.State.CurrentBranch,
			Remote:  remote,
		})
	}

	m.loadSelectedDetails()
}

func (m BranchesModel) Update(msg tea.Msg) (BranchesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "backspace":
			m.State.Done = true

		case "r":
			m.Load(m.State.RepoPath)

		case "j", "down":
			m.moveDown()

		case "k", "up":
			m.moveUp()

		case "tab":
			m.togglePanel()

		case "enter":
			m.checkoutSelected()

		case "m":
			m.mergeSelected(false)

		case "M":
			m.mergeSelected(true)

		case "d":
			m.deleteSelected()
		}
	}

	return m, nil
}

func (m BranchesModel) SelectedBranch() string {
	if len(m.State.Branches) == 0 ||
		m.State.Selected < 0 ||
		m.State.Selected >= len(m.State.Branches) {
		return ""
	}

	return m.State.Branches[m.State.Selected].Name
}

func (m *BranchesModel) moveDown() {
	if m.State.Selected < len(m.State.Branches)-1 {
		m.State.Selected++
		m.loadSelectedDetails()
	}
}

func (m *BranchesModel) moveUp() {
	if m.State.Selected > 0 {
		m.State.Selected--
		m.loadSelectedDetails()
	}
}

func (m *BranchesModel) togglePanel() {
	if m.State.Panel == 0 {
		m.State.Panel = 1
	} else {
		m.State.Panel = 0
	}
}

func (m *BranchesModel) checkoutSelected() {
	branch, ok := m.selectedBranch()
	if !ok {
		return
	}

	if branch.Name == m.State.CurrentBranch {
		return
	}

	out, err := m.Git.GitOutput(m.State.RepoPath, "checkout", branch.Name)
	if err != nil {
		m.State.Error = cleanGitError(out, err)
		return
	}

	m.State.CheckedOut = true
	m.reloadPreservingActionFlags()
}

func (m *BranchesModel) mergeSelected(deleteAfter bool) {
	branch, ok := m.selectedBranch()
	if !ok {
		return
	}

	if branch.Name == m.State.CurrentBranch {
		m.State.Error = "cannot merge current branch into itself"
		return
	}

	if branch.Remote {
		m.State.Error = "remote-tracking branch merge is disabled for now"
		return
	}

	out, err := m.Git.GitOutput(m.State.RepoPath, "merge", "--no-ff", branch.Name)
	if err != nil {
		m.State.Error = cleanGitError(out, err)
		return
	}

	m.State.Merged = true

	if deleteAfter {
		out, err := m.Git.GitOutput(m.State.RepoPath, "branch", "-d", branch.Name)
		if err != nil {
			m.State.Error = cleanGitError(out, err)
			return
		}

		m.State.Deleted = true
	}

	m.reloadPreservingActionFlags()
}

func (m *BranchesModel) deleteSelected() {
	branch, ok := m.selectedBranch()
	if !ok {
		return
	}

	if branch.Name == m.State.CurrentBranch {
		m.State.Error = "cannot delete current branch"
		return
	}

	if branch.Remote {
		m.State.Error = "remote-tracking branch delete is disabled for now"
		return
	}

	out, err := m.Git.GitOutput(m.State.RepoPath, "branch", "-d", branch.Name)
	if err != nil {
		m.State.Error = cleanGitError(out, err)
		return
	}

	m.State.Deleted = true
	m.reloadPreservingActionFlags()
}

func (m *BranchesModel) loadSelectedDetails() {
	if len(m.State.Branches) == 0 {
		m.State.Commits = nil
		m.State.Files = nil
		return
	}

	branch := m.State.Branches[m.State.Selected].Name

	commits, err := m.Git.GitOutput(
		m.State.RepoPath,
		"log",
		branch,
		"--oneline",
		"--decorate",
		"--graph",
		"-n",
		"25",
	)
	if err != nil {
		m.State.Commits = []string{err.Error()}
	} else {
		m.State.Commits = cleanBranchLines(commits)
	}

	files, err := m.Git.GitOutput(
		m.State.RepoPath,
		"ls-tree",
		"-r",
		"--name-only",
		branch,
	)
	if err != nil {
		m.State.Files = []string{err.Error()}
	} else {
		m.State.Files = cleanBranchLinesLimit(files, 80)
	}
}

func (m *BranchesModel) reloadPreservingActionFlags() {
	merged := m.State.Merged
	deleted := m.State.Deleted
	checkedOut := m.State.CheckedOut
	repoPath := m.State.RepoPath

	m.Load(repoPath)

	m.State.Merged = merged
	m.State.Deleted = deleted
	m.State.CheckedOut = checkedOut
}

func (m BranchesModel) selectedBranch() (models.GitBranch, bool) {
	if len(m.State.Branches) == 0 ||
		m.State.Selected < 0 ||
		m.State.Selected >= len(m.State.Branches) {
		return models.GitBranch{}, false
	}

	return m.State.Branches[m.State.Selected], true
}

func containsBranch(branches []models.GitBranch, name string) bool {
	for _, branch := range branches {
		if branch.Name == name {
			return true
		}
	}

	return false
}

func cleanGitError(out string, err error) string {
	message := strings.TrimSpace(out)
	if message == "" && err != nil {
		message = err.Error()
	}

	return message
}

func cleanBranchLines(s string) []string {
	var out []string

	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimRight(line, "\r\n")
		if strings.TrimSpace(line) != "" {
			out = append(out, line)
		}
	}

	return out
}

func cleanBranchLinesLimit(s string, limit int) []string {
	lines := cleanBranchLines(s)

	if len(lines) > limit {
		return lines[:limit]
	}

	return lines
}
