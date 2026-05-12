package controllers

import (
	tea "github.com/charmbracelet/bubbletea"

	"ingit/enums"
	"ingit/interfaces"
	"ingit/models"
)

func (m Model) loadDiffCmd() tea.Cmd {
	if len(m.State.Repos) == 0 || m.State.ActiveRepo < 0 || m.State.ActiveRepo >= len(m.State.Repos) {
		return nil
	}

	return func() tea.Msg {
		repo := m.State.Repos[m.State.ActiveRepo]
		state := m.currentState()

		file, ok := m.selectedFile()
		if ok {
			state.Diff = m.Git.ReadDiffForFile(repo, file)
		} else {
			state.Diff = ""
		}

		return models.RepoLoadedMsg{Index: m.State.ActiveRepo, State: state}
	}
}

func loadRepoCmd(git interfaces.GitProvider, index int, repo models.Repo) tea.Cmd {
	return func() tea.Msg {
		state := git.ReadGitState(repo)

		if len(state.Files) > 0 {
			state.Diff = git.ReadDiffForFile(repo, state.Files[0])
		}

		return models.RepoLoadedMsg{Index: index, State: state}
	}
}

func (m Model) executeQuickPushCmd(message string) tea.Cmd {
	repo := m.State.Repos[m.State.ActiveRepo]

	return func() tea.Msg {
		err := m.Git.ExecuteQuickPush(repo, message)
		return models.QuickPushExecutedMsg{Err: err}
	}
}

func (m Model) executeSelectedPlanCmd() tea.Cmd {
	if m.State.Mode != enums.ModePush {
		return nil
	}

	if len(m.State.Plan.Branches) == 0 {
		return nil
	}

	if m.State.ActiveBranch < 0 || m.State.ActiveBranch >= len(m.State.Plan.Branches) {
		return nil
	}

	index := m.State.ActiveBranch
	plan := m.State.Plan.Branches[index]
	repo := m.State.Repos[m.State.ActiveRepo]

	return func() tea.Msg {
		err := m.Git.ExecutePushPlan(repo, models.PushPlanExecution{
			Branch:  plan.Name,
			Base:    plan.Base,
			Message: plan.Message,
			Files:   plan.Files,
		})

		return models.PlanExecutedMsg{Index: index, Err: err}
	}
}

func (m Model) executeAllPlansCmd() tea.Cmd {
	if m.State.Mode != enums.ModePush {
		return nil
	}

	if len(m.State.Plan.Branches) == 0 {
		return nil
	}

	plans := append([]interfaces.PlannedBranch(nil), m.State.Plan.Branches...)
	repo := m.State.Repos[m.State.ActiveRepo]

	return func() tea.Msg {
		var results []models.PlanExecutedMsg

		for index, plan := range plans {
			err := m.Git.ExecutePushPlan(repo, models.PushPlanExecution{
				Branch:  plan.Name,
				Base:    plan.Base,
				Message: plan.Message,
				Files:   plan.Files,
			})

			results = append(results, models.PlanExecutedMsg{
				Index: index,
				Err:   err,
			})

			if err != nil {
				break
			}
		}

		return models.AllPlansExecutedMsg{Results: results}
	}
}

func (m Model) stageSelectedCmd() tea.Cmd {
	file, ok := m.selectedFile()
	if !ok {
		return nil
	}

	repo := m.State.Repos[m.State.ActiveRepo]

	return runActionCmd(func() error {
		return m.Git.StageFile(repo, file.Path)
	}, "staged "+file.Path)
}

func (m Model) unstageSelectedCmd() tea.Cmd {
	file, ok := m.selectedFile()
	if !ok {
		return nil
	}

	repo := m.State.Repos[m.State.ActiveRepo]

	return runActionCmd(func() error {
		return m.Git.UnstageFile(repo, file.Path)
	}, "unstaged "+file.Path)
}

func (m Model) stageAllCmd() tea.Cmd {
	repo := m.State.Repos[m.State.ActiveRepo]

	return runActionCmd(func() error {
		return m.Git.StageAll(repo)
	}, "staged all")
}

func (m Model) unstageAllCmd() tea.Cmd {
	repo := m.State.Repos[m.State.ActiveRepo]

	return runActionCmd(func() error {
		return m.Git.UnstageAll(repo)
	}, "unstaged all")
}

func runActionCmd(fn func() error, message string) tea.Cmd {
	return func() tea.Msg {
		err := fn()
		return models.ActionDoneMsg{
			Message: message,
			Err:     err,
		}
	}
}
