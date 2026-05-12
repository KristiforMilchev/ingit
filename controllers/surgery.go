package controllers

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"ingit/interfaces"
	"ingit/models"
	"ingit/states"
)

type SurgeryModel struct {
	State states.SurgeryState
	Git   interfaces.GitProvider
}

func NewSurgeryModel(git interfaces.GitProvider) SurgeryModel {
	return SurgeryModel{
		State: states.SurgeryState{},
		Git:   git,
	}
}

func (m *SurgeryModel) Load(repoPath string, filePath string) {
	m.State = states.SurgeryState{
		RepoPath: repoPath,
		FilePath: filePath,
	}

	diff, err := m.Git.GitOutput(repoPath, "diff", "--", filePath)
	if err != nil {
		m.State.Error = cleanGitError(diff, err)
		return
	}

	if strings.TrimSpace(diff) == "" {
		diff, err = m.Git.GitOutput(repoPath, "diff", "--cached", "--", filePath)
		if err != nil {
			m.State.Error = cleanGitError(diff, err)
			return
		}
	}

	m.State.Hunks = parseDiffHunks(filePath, diff)
}

func (m SurgeryModel) Update(msg tea.Msg) (SurgeryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.State.WritingMessage {
			return m.handleMessageInput(msg), nil
		}

		switch msg.String() {
		case "esc", "backspace":
			m.State.Done = true

		case "j", "down":
			if m.State.Selected < len(m.State.Hunks)-1 {
				m.State.Selected++
			}

		case "k", "up":
			if m.State.Selected > 0 {
				m.State.Selected--
			}

		case " ":
			m.toggleSelectedHunk()

		case "C":
			if m.selectedHunkCount() == 0 {
				m.State.Error = "select at least one hunk"
				return m, nil
			}
			m.State.WritingMessage = true
			m.State.Error = ""

		case "e":
			if err := m.executeCommitSlice(); err != nil {
				m.State.Error = err.Error()
				return m, nil
			}
			m.State.Executed = true
		}
	}

	return m, nil
}

func (m SurgeryModel) handleMessageInput(msg tea.KeyMsg) SurgeryModel {
	switch msg.String() {
	case "esc":
		m.State.WritingMessage = false
		m.State.MessageInput = ""
		return m

	case "enter":
		message := strings.TrimSpace(m.State.MessageInput)
		if message == "" {
			m.State.Error = "commit message cannot be empty"
			return m
		}

		if err := m.executeCommitSlice(); err != nil {
			m.State.Error = err.Error()
			return m
		}

		m.State.Executed = true
		return m

	case "backspace":
		if len(m.State.MessageInput) > 0 {
			runes := []rune(m.State.MessageInput)
			m.State.MessageInput = string(runes[:len(runes)-1])
		}
		return m

	default:
		if len(msg.Runes) > 0 {
			m.State.MessageInput += string(msg.Runes)
		}
		return m
	}
}

func (m *SurgeryModel) toggleSelectedHunk() {
	if len(m.State.Hunks) == 0 {
		return
	}

	m.State.Hunks[m.State.Selected].Selected = !m.State.Hunks[m.State.Selected].Selected
	m.State.Error = ""
}

func (m SurgeryModel) selectedHunkCount() int {
	count := 0
	for _, hunk := range m.State.Hunks {
		if hunk.Selected {
			count++
		}
	}
	return count
}

func parseDiffHunks(filePath string, diff string) []models.DiffHunk {
	var hunks []models.DiffHunk
	var current *models.DiffHunk

	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "@@") {
			if current != nil {
				hunks = append(hunks, *current)
			}

			current = &models.DiffHunk{
				FilePath: filePath,
				Header:   line,
				Lines:    []string{},
			}
			continue
		}

		if current != nil {
			current.Lines = append(current.Lines, line)
		}
	}

	if current != nil {
		hunks = append(hunks, *current)
	}

	return hunks
}

func (m SurgeryModel) executeCommitSlice() error {
	message := strings.TrimSpace(m.State.MessageInput)
	if message == "" {
		return fmt.Errorf("commit message cannot be empty")
	}

	patch := m.buildSelectedPatch()
	if strings.TrimSpace(patch) == "" {
		return fmt.Errorf("no hunks selected")
	}

	_, err := m.Git.GitOutputWithInput(m.State.RepoPath, patch, "apply", "--cached")
	if err != nil {
		return err
	}

	out, err := m.Git.GitOutput(m.State.RepoPath, "commit", "-m", message)
	if err != nil {
		return fmt.Errorf(cleanGitError(out, err))
	}

	return nil
}

func (m SurgeryModel) buildSelectedPatch() string {
	var b strings.Builder

	b.WriteString("diff --git a/")
	b.WriteString(m.State.FilePath)
	b.WriteString(" b/")
	b.WriteString(m.State.FilePath)
	b.WriteString("\n")

	b.WriteString("--- a/")
	b.WriteString(m.State.FilePath)
	b.WriteString("\n")

	b.WriteString("+++ b/")
	b.WriteString(m.State.FilePath)
	b.WriteString("\n")

	for _, hunk := range m.State.Hunks {
		if !hunk.Selected {
			continue
		}

		b.WriteString(hunk.Header)
		b.WriteString("\n")

		for _, line := range hunk.Lines {
			b.WriteString(line)
			b.WriteString("\n")
		}
	}

	return b.String()
}
