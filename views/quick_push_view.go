package views

import (
	"fmt"
	"strings"

	"ingit/states"
	"ingit/utils"
)

func renderQuickPushPanel(state states.AppState, w, h int) string {
	repoState := currentState(state)

	lines := []string{
		TitleStyle.Render("Quick Push"),
		MutedStyle.Render("current branch: " + repoState.Branch),
		MutedStyle.Render("operation: git add -A && git commit && git push"),
		"",
		TitleStyle.Render("Commit message"),
		InputStyle.Render(utils.TruncateLine(" "+state.QuickMessageInput+" ", w-4)),
		MutedStyle.Render("enter push · esc cancel"),
		"",
		TitleStyle.Render("Files included"),
	}

	if len(repoState.Files) == 0 {
		lines = append(lines, MutedStyle.Render("clean"))
	} else {
		for _, file := range repoState.Files {
			status := fmt.Sprintf("%s%s", normalizeStatus(file.Index), normalizeStatus(file.Worktree))
			line := fmt.Sprintf("%-2s %s", status, file.Path)
			lines = append(lines, OkStyle.Render(utils.TruncateLine("  "+line, w-4)))
		}
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)

	return FocusedBorder.Width(w).Height(h).Render(content)
}
