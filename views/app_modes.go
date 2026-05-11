package views

import (
	"github.com/charmbracelet/lipgloss"

	"ingit/states"
)

func renderQuickPushMode(state states.AppState, w, h int) string {
	leftW := 54
	if w >= 170 {
		leftW = 66
	}

	rightW := w - leftW - 2
	if rightW < 60 {
		rightW = 60
	}

	left := renderQuickPushPanel(state, leftW, h)
	right := renderRepositoryGraph(state, rightW, h)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func renderGitMode(state states.AppState, w, h int) string {
	if len(state.Repos) == 1 {
		return renderSingleRepoView(state, w, h)
	}

	leftW := responsiveLeftWidth(w)
	rightW := w - leftW - 2
	if rightW < 60 {
		rightW = 60
	}

	repoH := h / 2
	filesH := h - repoH

	left := lipgloss.JoinVertical(
		lipgloss.Left,
		renderRepos(state, leftW, repoH),
		renderFiles(state, leftW, filesH),
	)

	right := renderDiff(state, rightW, h)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func renderSingleRepoView(state states.AppState, w, h int) string {
	filesW := responsiveLeftWidth(w)
	diffW := w - filesW - 2
	if diffW < 60 {
		diffW = 60
	}

	files := renderFiles(state, filesW, h)
	diff := renderDiff(state, diffW, h)

	return lipgloss.JoinHorizontal(lipgloss.Top, files, diff)
}

func renderPushMode(state states.AppState, w, h int) string {
	planW := 42
	fileW := 46

	if w >= 180 {
		planW = 52
		fileW = 58
	}

	graphW := w - planW - fileW - 3
	if graphW < 60 {
		graphW = 60
	}

	plan := renderPushPlans(state, planW, h)
	graph := renderRepositoryGraph(state, graphW, h)
	files := renderAssignmentFiles(state, fileW, h)

	return lipgloss.JoinHorizontal(lipgloss.Top, plan, graph, files)
}
