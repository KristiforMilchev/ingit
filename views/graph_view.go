package views

import (
	"fmt"
	"strings"

	"ingit/enums"
	"ingit/states"
	"ingit/utils"
)

func renderRepositoryGraph(state states.AppState, w, h int) string {
	mode := "git graph"
	if state.GraphMode == enums.GraphText {
		mode = "branch text"
	}

	lines := []string{
		TitleStyle.Render("Repository Branches") + "  " + MutedStyle.Render("("+mode+", g toggles)"),
	}

	if state.GraphMode == enums.GraphText {
		lines = append(lines, renderBranchTextLines(state, w-4)...)
	} else {
		graph := currentState(state).Graph
		if strings.TrimSpace(graph) == "" {
			lines = append(lines, MutedStyle.Render("no graph available"))
		} else {
			for _, line := range strings.Split(graph, "\n") {
				lines = append(lines, renderGraphLine(utils.TruncateLine(line, w-4)))
			}
		}
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)

	style := PanelBorder.Width(w).Height(h)
	if state.Focus == enums.PanelGraph {
		style = FocusedBorder.Width(w).Height(h)
	}

	return style.Render(content)
}

func renderBranchTextLines(state states.AppState, width int) []string {
	repoState := currentState(state)

	current := repoState.Branch
	if current == "" {
		current = "current"
	}

	lines := []string{
		OkStyle.Render(utils.TruncateLine("current: "+current, width)),
		MutedStyle.Render(""),
		TitleStyle.Render("planned push branches"),
	}

	if len(state.Plan.Branches) == 0 {
		lines = append(lines, MutedStyle.Render("  none"))
		return lines
	}

	for i, branch := range state.Plan.Branches {
		style := planStyle(i)

		lines = append(lines, style.Render(utils.TruncateLine("  ├─ "+branch.Name, width)))
		lines = append(lines, MutedStyle.Render(utils.TruncateLine("  │  base: "+branch.Base, width)))

		filesText := fmt.Sprintf("  │  files: %d", len(branch.Files))
		lines = append(lines, MutedStyle.Render(utils.TruncateLine(filesText, width)))

		message := branch.Message
		if strings.TrimSpace(message) == "" {
			message = "<missing message>"
		}

		lines = append(lines, MutedStyle.Render(utils.TruncateLine("  │  msg: "+message, width)))
	}

	return lines
}

func renderGraphLine(line string) string {
	switch {
	case strings.Contains(line, "HEAD"):
		return OkStyle.Render(line)
	case strings.Contains(line, "origin/"):
		return WarnStyle.Render(line)
	case strings.Contains(line, "tag:"):
		return TitleStyle.Render(line)
	default:
		return MutedStyle.Render(line)
	}
}
