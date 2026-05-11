package views

import (
	"github.com/charmbracelet/lipgloss"

	"ingit/enums"
	"ingit/states"
)

func RenderAppView(state states.AppState, branches states.BranchesState) string {
	if state.Width == 0 || state.Height == 0 {
		return "loading..."
	}

	if len(state.Repos) == 0 {
		return "no repositories configured"
	}

	topH := state.Height - 6
	if topH < 10 {
		topH = 10
	}

	header := renderHeader(state)
	var body string

	switch state.Mode {
	case enums.ModePush:
		body = renderPushMode(state, state.Width, topH)
	case enums.ModeQuickPush:
		body = renderQuickPushMode(state, state.Width, topH)
	case enums.ModeBranches:
		body = RenderBranchesView(branches, state.Width, topH)
	default:
		body = renderGitMode(state, state.Width, topH)
	}

	footer := renderFooter(state, state.Width)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}
