package views

import (
	"fmt"

	"ingit/enums"
	"ingit/states"
	"ingit/utils"
)

func renderHeader(state states.AppState) string {
	repo := state.Repos[state.ActiveRepo]
	repoState := currentState(state)

	branch := repoState.Branch
	if branch == "" {
		branch = "no branch"
	}

	mode := "git"
	if state.Mode == enums.ModePush {
		mode = "push-plan"
	}
	if state.Mode == enums.ModeQuickPush {
		mode = "quick-push"
	}
	if state.Mode == enums.ModeBranches {
		mode = "branches"
	}
	if state.Mode == enums.ModeSurgery {
		mode = "surgery"
	}

	width := state.Width
	if width < 20 {
		width = 20
	}

	textWidth := width - 4
	if textWidth < 10 {
		textWidth = 10
	}

	text := fmt.Sprintf(
		" INGIT │ %s │ %s │ %s │ files:%d ",
		mode,
		repo.Name,
		branch,
		len(repoState.Files),
	)

	text = utils.TruncateLine(text, textWidth)

	return StatusBarStyle.
		Width(width).
		MaxWidth(width).
		Render(text)
}

func renderFooter(state states.AppState, w int) string {
	if w < 20 {
		w = 20
	}

	var commands string

	switch state.Mode {
	case enums.ModeQuickPush:
		commands = "enter commit+push · esc cancel · type message"
	case enums.ModeBranches:
		commands = "enter checkout · m merge · M merge+delete · d delete · tab commits/files · j/k move · r refresh · esc back · q"
	case enums.ModeSurgery:
		commands = "space select hunk · C message · e execute · j/k move · esc back"
	case enums.ModePush:
		commands = "esc git · g graph/text · C msg · e push · E all · tab panel · space mark · b plan · a assign · x unassign · r refresh · q"
	default:
		commands = "p planner · P quick push · B branches · I surgery · j/k move · h/l/tab panel · s/u stage · S/U all · r refresh · q"
	}

	status := state.StatusMsg

	if state.NamingBranch {
		if state.QuickPush {
			status = "quick push branch: " + state.BranchInput
		} else {
			status = "branch: " + state.BranchInput
		}
	}

	if state.NamingMessage {
		status = "message: " + state.MessageInput
	}

	if state.Mode == enums.ModeQuickPush {
		status = "quick message: " + state.QuickMessageInput
	}

	if status == "" {
		status = "ready"
	}

	lineWidth := w - 6
	if lineWidth < 10 {
		lineWidth = 10
	}

	commands = utils.TruncateLine(commands, lineWidth)
	status = utils.TruncateLine(status, lineWidth)

	return PanelBorder.
		Width(w - 2).
		MaxWidth(w - 2).
		Render(commands + "\n" + MutedStyle.Render(status))
}

func responsiveLeftWidth(w int) int {
	leftW := 38

	if w >= 150 {
		leftW = 44
	}

	if w >= 190 {
		leftW = 52
	}

	return leftW
}
