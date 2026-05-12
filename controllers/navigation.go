package controllers

import "ingit/enums"

func (m *Model) nextPanel() {
	if m.State.Mode == enums.ModePush {
		switch m.State.Focus {
		case enums.PanelPlans:
			m.State.Focus = enums.PanelGraph
		case enums.PanelGraph:
			m.State.Focus = enums.PanelFiles
		default:
			m.State.Focus = enums.PanelPlans
		}
		return
	}

	if len(m.State.Repos) == 1 {
		m.State.Focus = enums.PanelFiles
		return
	}

	if m.State.Focus == enums.PanelRepos {
		m.State.Focus = enums.PanelFiles
	} else {
		m.State.Focus = enums.PanelRepos
	}
}

func (m *Model) prevPanel() {
	if m.State.Mode == enums.ModePush {
		switch m.State.Focus {
		case enums.PanelPlans:
			m.State.Focus = enums.PanelFiles
		case enums.PanelGraph:
			m.State.Focus = enums.PanelPlans
		default:
			m.State.Focus = enums.PanelGraph
		}
		return
	}

	m.nextPanel()
}

func (m *Model) moveDown() {
	switch m.State.Focus {
	case enums.PanelRepos:
		if m.State.ActiveRepo < len(m.State.Repos)-1 {
			m.State.ActiveRepo++
			m.State.ActiveFile = 0
		}
	case enums.PanelFiles:
		files := m.currentState().Files
		if m.State.ActiveFile < len(files)-1 {
			m.State.ActiveFile++
		}
	case enums.PanelPlans:
		if m.State.ActiveBranch < len(m.State.Plan.Branches)-1 {
			m.State.ActiveBranch++
		}
	}
}

func (m *Model) moveUp() {
	switch m.State.Focus {
	case enums.PanelRepos:
		if m.State.ActiveRepo > 0 {
			m.State.ActiveRepo--
			m.State.ActiveFile = 0
		}
	case enums.PanelFiles:
		if m.State.ActiveFile > 0 {
			m.State.ActiveFile--
		}
	case enums.PanelPlans:
		if m.State.ActiveBranch > 0 {
			m.State.ActiveBranch--
		}
	}
}

func (m *Model) startQuickPushView() {
	m.State.Mode = enums.ModeQuickPush
	m.State.Focus = enums.PanelFiles
	m.State.QuickMessageInput = ""
	m.State.StatusMsg = "quick push: type commit message, enter to commit+push, esc to cancel"
}

func (m *Model) startBranchNameInput() {
	m.State.Mode = enums.ModePush
	m.State.Focus = enums.PanelPlans
	m.State.NamingBranch = true
	m.State.QuickPush = false
	m.State.BranchInput = "feature/"
	m.State.StatusMsg = "enter push plan branch name"
}

func (m *Model) startQuickPushInput() {
	m.State.Mode = enums.ModePush
	m.State.Focus = enums.PanelPlans
	m.State.NamingBranch = true
	m.State.QuickPush = true
	m.State.BranchInput = "feature/"
	m.State.StatusMsg = "quick push: name branch, all changed files will be assigned"
}

func (m *Model) startMessageInput() {
	m.State.Mode = enums.ModePush
	m.State.Focus = enums.PanelPlans

	if len(m.State.Plan.Branches) == 0 {
		m.State.StatusMsg = "create a push plan first"
		return
	}

	if m.State.ActiveBranch < 0 || m.State.ActiveBranch >= len(m.State.Plan.Branches) {
		m.State.ActiveBranch = 0
	}

	m.State.NamingMessage = true
	m.State.MessageInput = m.State.Plan.Branches[m.State.ActiveBranch].Message
	m.State.StatusMsg = "enter commit message for " + m.State.Plan.Branches[m.State.ActiveBranch].Name
}
