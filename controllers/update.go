package controllers

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"ingit/enums"
	"ingit/interfaces"
	"ingit/models"
)

func (m Model) Init() tea.Cmd {
	if len(m.State.Repos) == 0 {
		return tickCmd()
	}

	return tea.Batch(loadRepoCmd(m.Git, 0, m.State.Repos[0]), tickCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.State.Width = msg.Width
		m.State.Height = msg.Height
		return m, nil

	case models.TickMsg:
		if len(m.State.Repos) == 0 {
			return m, tickCmd()
		}

		if m.State.Focus == enums.PanelFiles {
			return m, tickCmd()
		}

		return m, tea.Batch(
			loadRepoCmd(m.Git, m.State.ActiveRepo, m.State.Repos[m.State.ActiveRepo]),
			tickCmd(),
		)

	case models.RepoLoadedMsg:
		if msg.Index < 0 || msg.Index >= len(m.State.Repos) {
			return m, nil
		}

		repo := m.State.Repos[msg.Index]

		m.State.States[repo.Path] = msg.State

		files := m.State.States[repo.Path].Files

		if m.State.ActiveFile >= len(files) {
			m.State.ActiveFile = len(files) - 1
		}

		if m.State.ActiveFile < 0 {
			m.State.ActiveFile = 0
		}

		if m.State.StatusMsg == "" || strings.HasPrefix(m.State.StatusMsg, "refreshed") {
			m.State.StatusMsg = "refreshed " + repo.Name
		}

		return m, nil

	case models.ActionDoneMsg:
		if len(m.State.Repos) == 0 {
			return m, nil
		}

		if msg.Err != nil {
			m.State.StatusMsg = msg.Err.Error()
		} else {
			m.State.StatusMsg = msg.Message
		}

		return m, loadRepoCmd(m.Git, m.State.ActiveRepo, m.State.Repos[m.State.ActiveRepo])

	case models.PlanExecutedMsg:
		if len(m.State.Repos) == 0 {
			return m, nil
		}

		if msg.Index >= 0 && msg.Index < len(m.State.Plan.Branches) {
			if msg.Err != nil {
				m.State.Plan.Branches[msg.Index].Status = "failed"
				m.State.Plan.Branches[msg.Index].Error = msg.Err.Error()
				m.State.StatusMsg = msg.Err.Error()
			} else {
				executed := m.State.Plan.Branches[msg.Index].Name
				m.State.Plan.Branches = append(m.State.Plan.Branches[:msg.Index], m.State.Plan.Branches[msg.Index+1:]...)
				m.State.MarkedFiles = map[string]bool{}

				if len(m.State.Plan.Branches) == 0 {
					m.State.ActiveBranch = 0
				} else if m.State.ActiveBranch >= len(m.State.Plan.Branches) {
					m.State.ActiveBranch = len(m.State.Plan.Branches) - 1
				}

				m.State.StatusMsg = "executed and cleared " + executed
			}
		}

		return m, loadRepoCmd(m.Git, m.State.ActiveRepo, m.State.Repos[m.State.ActiveRepo])

	case models.QuickPushExecutedMsg:
		if len(m.State.Repos) == 0 {
			return m, nil
		}

		if msg.Err != nil {
			m.State.StatusMsg = msg.Err.Error()
			return m, loadRepoCmd(m.Git, m.State.ActiveRepo, m.State.Repos[m.State.ActiveRepo])
		}

		m.State.Mode = enums.ModeGit
		m.State.Focus = enums.PanelFiles
		m.State.QuickMessageInput = ""
		m.State.MarkedFiles = map[string]bool{}
		m.State.StatusMsg = "quick push completed"

		return m, loadRepoCmd(m.Git, m.State.ActiveRepo, m.State.Repos[m.State.ActiveRepo])

	case models.AllPlansExecutedMsg:
		if len(m.State.Repos) == 0 {
			return m, nil
		}

		failed := 0
		failedByIndex := map[int]error{}

		for _, result := range msg.Results {
			if result.Err != nil {
				failed++
				failedByIndex[result.Index] = result.Err
			}
		}

		var remaining []interfaces.PlannedBranch
		for index, branch := range m.State.Plan.Branches {
			if err, ok := failedByIndex[index]; ok {
				branch.Status = "failed"
				branch.Error = err.Error()
				remaining = append(remaining, branch)
			}
		}

		m.State.Plan.Branches = remaining
		m.State.MarkedFiles = map[string]bool{}

		if len(m.State.Plan.Branches) == 0 {
			m.State.ActiveBranch = 0
		} else if m.State.ActiveBranch >= len(m.State.Plan.Branches) {
			m.State.ActiveBranch = len(m.State.Plan.Branches) - 1
		}

		if failed > 0 {
			m.State.StatusMsg = fmt.Sprintf("executed all; %d failed plan(s) kept", failed)
		} else {
			m.State.StatusMsg = "all push plans executed and cleared"
		}

		return m, loadRepoCmd(m.Git, m.State.ActiveRepo, m.State.Repos[m.State.ActiveRepo])

	case tea.KeyMsg:
		if len(m.State.Repos) == 0 {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			}
			return m, nil
		}

		if m.State.Mode == enums.ModeSurgery {
			var cmd tea.Cmd
			m.Surgery, cmd = m.Surgery.Update(msg)

			if m.Surgery.State.Done {
				m.Surgery.State.Done = false
				m.State.Mode = enums.ModeGit
				m.State.Focus = enums.PanelFiles
				m.State.StatusMsg = "git view"
				return m, loadRepoCmd(m.Git, m.State.ActiveRepo, m.State.Repos[m.State.ActiveRepo])
			}

			if m.Surgery.State.Executed {
				m.Surgery.State.Executed = false
				m.State.Mode = enums.ModeGit
				m.State.Focus = enums.PanelFiles
				m.State.StatusMsg = "commit surgery executed"
				return m, loadRepoCmd(m.Git, m.State.ActiveRepo, m.State.Repos[m.State.ActiveRepo])
			}

			return m, cmd
		}

		if m.State.Mode == enums.ModeBranches {
			var cmd tea.Cmd
			m.Branches, cmd = m.Branches.Update(msg)

			if m.Branches.State.Done {
				m.Branches.State.Done = false
				m.State.Mode = enums.ModeGit
				m.State.Focus = enums.PanelFiles
				m.State.StatusMsg = "git view"
				return m, nil
			}

			if m.Branches.State.CheckedOut {
				branch := m.Branches.SelectedBranch()
				m.Branches.State.CheckedOut = false
				m.State.Mode = enums.ModeGit
				m.State.Focus = enums.PanelFiles

				if branch == "" {
					m.State.StatusMsg = "branch switched"
				} else {
					m.State.StatusMsg = "switched to " + branch
				}

				return m, loadRepoCmd(m.Git, m.State.ActiveRepo, m.State.Repos[m.State.ActiveRepo])
			}

			if m.Branches.State.Merged {
				m.Branches.State.Merged = false
				m.State.Mode = enums.ModeGit
				m.State.Focus = enums.PanelFiles

				if m.Branches.State.Deleted {
					m.Branches.State.Deleted = false
					m.State.StatusMsg = "merged and deleted branch"
				} else {
					m.State.StatusMsg = "merged branch"
				}

				return m, loadRepoCmd(m.Git, m.State.ActiveRepo, m.State.Repos[m.State.ActiveRepo])
			}

			if m.Branches.State.Deleted {
				m.Branches.State.Deleted = false
				m.State.StatusMsg = "deleted branch"
				return m, loadRepoCmd(m.Git, m.State.ActiveRepo, m.State.Repos[m.State.ActiveRepo])
			}

			return m, cmd
		}

		if m.State.Mode == enums.ModeQuickPush {
			return m.handleQuickPushInput(msg)
		}

		if m.State.NamingBranch {
			return m.handleBranchNameInput(msg)
		}

		if m.State.NamingMessage {
			return m.handleMessageInput(msg)
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "p":
			m.State.Mode = enums.ModePush
			m.State.Focus = enums.PanelPlans
			m.State.StatusMsg = "push planner"
			return m, nil

		case "P":
			m.startQuickPushView()
			return m, nil

		case "B":
			repo := m.State.Repos[m.State.ActiveRepo]
			m.Branches.Load(repo.Path)
			m.State.Mode = enums.ModeBranches
			m.State.Focus = enums.PanelGraph
			m.State.StatusMsg = "branches"
			return m, nil

		case "I":
			file, ok := m.selectedFile()
			if !ok {
				m.State.StatusMsg = "select a file first"
				return m, nil
			}

			repo := m.State.Repos[m.State.ActiveRepo]
			m.Surgery.Load(repo.Path, file.Path)
			m.State.Mode = enums.ModeSurgery
			m.State.Focus = enums.PanelFiles
			m.State.StatusMsg = "commit surgery"
			return m, nil

		case "esc":
			m.State.Mode = enums.ModeGit
			m.State.Focus = enums.PanelFiles
			m.State.StatusMsg = "git view"
			return m, nil

		case "g":
			if m.State.GraphMode == enums.GraphVisual {
				m.State.GraphMode = enums.GraphText
				m.State.StatusMsg = "graph: branch text"
			} else {
				m.State.GraphMode = enums.GraphVisual
				m.State.StatusMsg = "graph: git log"
			}
			return m, nil

		case "tab", "l", "right":
			m.nextPanel()
			return m, nil

		case "h", "left":
			m.prevPanel()
			return m, nil

		case "j", "down":
			m.moveDown()
			return m, m.loadDiffCmd()

		case "k", "up":
			m.moveUp()
			return m, m.loadDiffCmd()

		case "r":
			m.State.StatusMsg = "refreshing..."
			return m, loadRepoCmd(m.Git, m.State.ActiveRepo, m.State.Repos[m.State.ActiveRepo])

		case " ":
			m.toggleMarkedFile()
			return m, nil

		case "b":
			m.startBranchNameInput()
			return m, nil

		case "C":
			m.startMessageInput()
			return m, nil

		case "a":
			m.assignMarkedToBranch()
			return m, nil

		case "x":
			m.removeSelectedFileFromPlan()
			return m, nil

		case "e":
			return m, m.executeSelectedPlanCmd()

		case "E":
			return m, m.executeAllPlansCmd()

		case "s":
			if m.State.Mode == enums.ModePush {
				m.State.StatusMsg = "staging is disabled in push planner"
				return m, nil
			}
			return m, m.stageSelectedCmd()

		case "u":
			if m.State.Mode == enums.ModePush {
				m.State.StatusMsg = "unstaging is disabled in push planner"
				return m, nil
			}
			return m, m.unstageSelectedCmd()

		case "S":
			if m.State.Mode == enums.ModePush {
				m.State.StatusMsg = "stage all is disabled in push planner"
				return m, nil
			}
			return m, m.stageAllCmd()

		case "U":
			if m.State.Mode == enums.ModePush {
				m.State.StatusMsg = "unstage all is disabled in push planner"
				return m, nil
			}
			return m, m.unstageAllCmd()
		}
	}

	return m, nil
}
