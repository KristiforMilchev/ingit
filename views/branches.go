package views

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ingit/utils"
)

type BranchesModel struct {
	RepoPath string

	Branches []GitBranch
	Selected int

	CurrentBranch string
	Panel         int

	Commits []string
	Files   []string

	Error      string
	Done       bool
	CheckedOut bool
}

type GitBranch struct {
	Name    string
	Current bool
	Remote  bool
}

func NewBranchesModel() BranchesModel {
	return BranchesModel{}
}

func (m *BranchesModel) Load(repoPath string) {
	m.RepoPath = repoPath
	m.Error = ""
	m.Done = false
	m.CheckedOut = false
	m.Branches = nil
	m.Commits = nil
	m.Files = nil
	m.Selected = 0

	current, _ := m.gitOutput("branch", "--show-current")
	m.CurrentBranch = strings.TrimSpace(current)

	out, err := m.gitOutput("branch", "--all", "--no-color")
	if err != nil {
		m.Error = strings.TrimSpace(out)
		if m.Error == "" {
			m.Error = err.Error()
		}
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

		if m.containsBranch(name) {
			continue
		}

		m.Branches = append(m.Branches, GitBranch{
			Name:    name,
			Current: current || name == m.CurrentBranch,
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
			m.Done = true

		case "r":
			m.Load(m.RepoPath)

		case "j", "down":
			if m.Selected < len(m.Branches)-1 {
				m.Selected++
				m.loadSelectedDetails()
			}

		case "k", "up":
			if m.Selected > 0 {
				m.Selected--
				m.loadSelectedDetails()
			}

		case "tab":
			if m.Panel == 0 {
				m.Panel = 1
			} else {
				m.Panel = 0
			}

		case "enter":
			m.checkoutSelected()
		}
	}

	return m, nil
}

func (m BranchesModel) View(width, height int) string {
	if m.Error != "" {
		content := TitleStyle.Render("Branches") +
			"\n\n" +
			ErrorStyle.Render(m.Error) +
			"\n\n" +
			MutedStyle.Render("esc/backspace back")

		return PanelBorder.Width(width - 2).Height(height).Render(content)
	}

	leftW := 44
	if width >= 170 {
		leftW = 56
	}

	rightW := width - leftW - 3
	if rightW < 60 {
		rightW = 60
	}

	left := m.renderBranchList(leftW, height)
	right := m.renderBranchDetails(rightW, height)

	return lipJoinHorizontal(left, right)
}

func (m BranchesModel) renderBranchList(w, h int) string {
	lines := []string{
		TitleStyle.Render("Branches"),
		MutedStyle.Render("enter checkout · tab details · esc back"),
		"",
	}

	if len(m.Branches) == 0 {
		lines = append(lines, MutedStyle.Render("no branches"))
	} else {
		for i, branch := range m.Branches {
			current := " "
			if branch.Current {
				current = "*"
			}

			remote := ""
			if branch.Remote {
				remote = " remote"
			}

			line := fmt.Sprintf("%s %s%s", current, branch.Name, remote)
			line = utils.TruncateLine(line, w-4)

			if i == m.Selected {
				line = SelectedStyle.Width(w - 4).Render(line)
			} else if branch.Current {
				line = OkStyle.Render(line)
			} else if branch.Remote {
				line = WarnStyle.Render(line)
			} else {
				line = MutedStyle.Render(line)
			}

			lines = append(lines, line)
		}
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)
	return FocusedBorder.Width(w).Height(h).Render(content)
}

func (m BranchesModel) renderBranchDetails(w, h int) string {
	branch := m.SelectedBranch()
	if branch == "" {
		return PanelBorder.Width(w).Height(h).Render(MutedStyle.Render("no branch selected"))
	}

	title := "Commits"
	if m.Panel == 1 {
		title = "Files"
	}

	lines := []string{
		TitleStyle.Render(title) + "  " + MutedStyle.Render(branch),
		"",
	}

	if m.Panel == 0 {
		if len(m.Commits) == 0 {
			lines = append(lines, MutedStyle.Render("no commits"))
		} else {
			for _, commit := range m.Commits {
				lines = append(lines, renderGraphLine(utils.TruncateLine(commit, w-4)))
			}
		}
	} else {
		if len(m.Files) == 0 {
			lines = append(lines, MutedStyle.Render("no files"))
		} else {
			for _, file := range m.Files {
				lines = append(lines, MutedStyle.Render(utils.TruncateLine(file, w-4)))
			}
		}
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)
	return PanelBorder.Width(w).Height(h).Render(content)
}

func (m *BranchesModel) checkoutSelected() {
	if len(m.Branches) == 0 {
		return
	}

	branch := m.Branches[m.Selected].Name
	if branch == m.CurrentBranch {
		return
	}

	out, err := m.gitOutput("checkout", branch)
	if err != nil {
		m.Error = strings.TrimSpace(out)
		if m.Error == "" {
			m.Error = err.Error()
		}
		return
	}

	m.CheckedOut = true
	m.Load(m.RepoPath)
}

func (m *BranchesModel) loadSelectedDetails() {
	if len(m.Branches) == 0 {
		m.Commits = nil
		m.Files = nil
		return
	}

	branch := m.Branches[m.Selected].Name

	commits, err := m.gitOutput("log", branch, "--oneline", "--decorate", "--graph", "-n", "25")
	if err != nil {
		m.Commits = []string{err.Error()}
	} else {
		m.Commits = cleanBranchLines(commits)
	}

	files, err := m.gitOutput("ls-tree", "-r", "--name-only", branch)
	if err != nil {
		m.Files = []string{err.Error()}
	} else {
		m.Files = cleanBranchLinesLimit(files, 80)
	}
}

func (m BranchesModel) SelectedBranch() string {
	if len(m.Branches) == 0 || m.Selected < 0 || m.Selected >= len(m.Branches) {
		return ""
	}

	return m.Branches[m.Selected].Name
}

func (m BranchesModel) containsBranch(name string) bool {
	for _, branch := range m.Branches {
		if branch.Name == name {
			return true
		}
	}

	return false
}

func (m BranchesModel) gitOutput(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = m.RepoPath

	out, err := cmd.CombinedOutput()
	return string(out), err
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

func lipJoinHorizontal(left string, right string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}
