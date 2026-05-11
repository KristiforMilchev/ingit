package controllers

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"ingit/models"
)

func tickCmd() tea.Cmd {
	return tea.Tick(900*time.Millisecond, func(t time.Time) tea.Msg {
		return models.TickMsg(t)
	})
}

func (m Model) currentState() models.GitState {
	if len(m.State.Repos) == 0 {
		return models.GitState{}
	}

	if m.State.ActiveRepo < 0 || m.State.ActiveRepo >= len(m.State.Repos) {
		return models.GitState{}
	}

	return m.State.States[m.State.Repos[m.State.ActiveRepo].Path]
}

func (m Model) selectedFile() (models.FileStatus, bool) {
	files := m.currentState().Files

	if len(files) == 0 || m.State.ActiveFile < 0 || m.State.ActiveFile >= len(files) {
		return models.FileStatus{}, false
	}

	return files[m.State.ActiveFile], true
}

func (m Model) selectedDiffModeLabel() string {
	file, ok := m.selectedFile()
	if !ok {
		return "diff"
	}

	if hasWorktreeChange(file) {
		return "unstaged diff"
	}

	if hasStagedChange(file) {
		return "staged diff"
	}

	return "diff"
}

func hasStagedChange(file models.FileStatus) bool {
	return file.Index != "" && file.Index != " " && file.Index != "." && file.Index != "?"
}

func hasWorktreeChange(file models.FileStatus) bool {
	return file.Worktree != "" && file.Worktree != " " && file.Worktree != "."
}
