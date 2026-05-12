package controllers

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"ingit/enums"
	"ingit/interfaces"
)

func (m Model) handleBranchNameInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.State.NamingBranch = false
		m.State.QuickPush = false
		m.State.BranchInput = ""
		m.State.StatusMsg = "branch creation cancelled"
		return m, nil

	case "enter":
		name := strings.TrimSpace(m.State.BranchInput)
		if name == "" {
			m.State.StatusMsg = "branch name cannot be empty"
			return m, nil
		}

		branch := interfaces.PlannedBranch{
			Name:   name,
			Base:   m.currentState().Branch,
			Status: "pending",
		}

		if m.State.QuickPush {
			for _, file := range m.currentState().Files {
				branch.Files = append(branch.Files, file.Path)
				m.removeFileFromAllBranches(file.Path)
			}
		}

		m.State.Plan.Branches = append(m.State.Plan.Branches, branch)
		m.State.ActiveBranch = len(m.State.Plan.Branches) - 1
		m.State.NamingBranch = false
		m.State.QuickPush = false
		m.State.BranchInput = ""

		if len(branch.Files) > 0 {
			m.State.StatusMsg = fmt.Sprintf("quick push plan %s with %d file(s)", name, len(branch.Files))
		} else {
			m.State.StatusMsg = "created push plan " + name
		}

		return m, nil

	case "backspace":
		if len(m.State.BranchInput) > 0 {
			runes := []rune(m.State.BranchInput)
			m.State.BranchInput = string(runes[:len(runes)-1])
		}
		return m, nil

	default:
		if len(msg.Runes) > 0 {
			m.State.BranchInput += string(msg.Runes)
		}
		return m, nil
	}
}

func (m Model) handleMessageInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.State.NamingMessage = false
		m.State.MessageInput = ""
		m.State.StatusMsg = "commit message cancelled"
		return m, nil

	case "enter":
		if len(m.State.Plan.Branches) == 0 {
			m.State.NamingMessage = false
			m.State.MessageInput = ""
			m.State.StatusMsg = "no push plan selected"
			return m, nil
		}

		message := strings.TrimSpace(m.State.MessageInput)
		if message == "" {
			m.State.StatusMsg = "commit message cannot be empty"
			return m, nil
		}

		if m.State.ActiveBranch < 0 || m.State.ActiveBranch >= len(m.State.Plan.Branches) {
			m.State.ActiveBranch = 0
		}

		m.State.Plan.Branches[m.State.ActiveBranch].Message = message
		m.State.NamingMessage = false
		m.State.MessageInput = ""
		m.State.StatusMsg = "message set for " + m.State.Plan.Branches[m.State.ActiveBranch].Name

		return m, nil

	case "backspace":
		if len(m.State.MessageInput) > 0 {
			runes := []rune(m.State.MessageInput)
			m.State.MessageInput = string(runes[:len(runes)-1])
		}
		return m, nil

	default:
		if len(msg.Runes) > 0 {
			m.State.MessageInput += string(msg.Runes)
		}
		return m, nil
	}
}

func (m Model) handleQuickPushInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.State.Mode = enums.ModeGit
		m.State.Focus = enums.PanelFiles
		m.State.QuickMessageInput = ""
		m.State.StatusMsg = "quick push cancelled"
		return m, nil

	case "enter":
		message := strings.TrimSpace(m.State.QuickMessageInput)
		if message == "" {
			m.State.StatusMsg = "commit message cannot be empty"
			return m, nil
		}

		return m, m.executeQuickPushCmd(message)

	case "backspace":
		if len(m.State.QuickMessageInput) > 0 {
			runes := []rune(m.State.QuickMessageInput)
			m.State.QuickMessageInput = string(runes[:len(runes)-1])
		}
		return m, nil

	default:
		if len(msg.Runes) > 0 {
			m.State.QuickMessageInput += string(msg.Runes)
		}
		return m, nil
	}
}
