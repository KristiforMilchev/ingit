package controllers

import (
	"ingit/enums"
	"ingit/interfaces"
	"ingit/models"
	"ingit/states"
)

type RenderAppFunc func(state states.AppState, branches states.BranchesState, surgery states.SurgeryState) string

type Model struct {
	Git      interfaces.GitProvider
	State    states.AppState
	Branches BranchesModel
	Surgery  SurgeryModel
	Render   RenderAppFunc
}

func NewModel(git interfaces.GitProvider, repos []models.Repo, render RenderAppFunc) Model {
	return Model{
		Git:      git,
		Branches: NewBranchesModel(git),
		Surgery:  NewSurgeryModel(git),
		Render:   render,
		State: states.AppState{
			Repos:       repos,
			States:      map[string]models.GitState{},
			ActiveRepo:  0,
			ActiveFile:  0,
			Focus:       enums.PanelFiles,
			Mode:        enums.ModeGit,
			GraphMode:   enums.GraphVisual,
			StatusMsg:   "ready",
			MarkedFiles: map[string]bool{},
			Plan:        interfaces.ChangePlan{Branches: []interfaces.PlannedBranch{}},
		},
	}
}

func (m Model) View() string {
	if m.Render == nil {
		return "missing renderer"
	}

	return m.Render(m.State, m.Branches.State, m.Surgery.State)
}
