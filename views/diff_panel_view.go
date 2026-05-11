package views

import (
	"fmt"

	"ingit/states"
	"ingit/utils"
)

func renderDiff(state states.AppState, w, h int) string {
	repoState := currentState(state)
	repo := state.Repos[state.ActiveRepo]
	mode := selectedDiffModeLabel(state)

	header := TitleStyle.Render("Diff") +
		"  " +
		MutedStyle.Render(repo.Name) +
		"  " +
		MutedStyle.Render(repoState.Branch) +
		"  " +
		MutedStyle.Render(mode)

	if repoState.Ahead > 0 || repoState.Behind > 0 {
		header += fmt.Sprintf("  ↑%d ↓%d", repoState.Ahead, repoState.Behind)
	}

	diff := repoState.Diff
	if diff == "" {
		diff = MutedStyle.Render("No diff selected.")
	} else if IsRealDiff(diff) {
		columnWidth := (w - 12) / 2
		if columnWidth >= 30 {
			diff = RenderSideBySideDiff(diff, columnWidth)
		} else {
			diff = RenderUnifiedDiff(diff, w-4)
		}
	} else {
		diff = MutedStyle.Render(utils.TruncateLine(diff, w-4))
	}

	content := utils.Fit(header+"\n\n"+diff, h-2)

	return PanelBorder.Width(w).Height(h).Render(content)
}
