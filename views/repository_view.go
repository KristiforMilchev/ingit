package views

import (
	"fmt"
	"strings"

	"ingit/enums"
	"ingit/states"
	"ingit/utils"
)

func renderRepos(state states.AppState, w, h int) string {
	lines := []string{TitleStyle.Render("Repositories")}

	nameWidth := w - 12
	if nameWidth < 10 {
		nameWidth = 10
	}

	for i, repo := range state.Repos {
		repoState := state.States[repo.Path]

		marker := " "
		if len(repoState.Files) > 0 {
			marker = WarnStyle.Render("*")
		} else if repoState.Error == "" && repoState.Branch != "" {
			marker = OkStyle.Render("✓")
		}

		name := utils.TruncateLine(repo.Name, nameWidth)
		line := fmt.Sprintf("%s %-*s %3d", marker, nameWidth, name, len(repoState.Files))

		if i == state.ActiveRepo {
			line = SelectedStyle.Width(w - 4).Render(line)
		}

		lines = append(lines, line)
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)

	style := PanelBorder.Width(w).Height(h)
	if state.Focus == enums.PanelRepos {
		style = FocusedBorder.Width(w).Height(h)
	}

	return style.Render(content)
}

func renderFiles(state states.AppState, w, h int) string {
	repoState := currentState(state)
	lines := []string{TitleStyle.Render("Files")}

	if repoState.Error != "" {
		lines = append(lines, ErrorStyle.Render(repoState.Error))
	} else if len(repoState.Files) == 0 {
		lines = append(lines, MutedStyle.Render("clean"))
	} else {
		pathWidth := w - 12
		if pathWidth < 12 {
			pathWidth = 12
		}

		for i, file := range repoState.Files {
			status := fmt.Sprintf("%s%s", normalizeStatus(file.Index), normalizeStatus(file.Worktree))

			mark := " "
			if state.MarkedFiles[file.Path] {
				mark = "●"
			}

			path := utils.TruncateLine(file.Path, pathWidth)
			line := fmt.Sprintf("%s %-2s %-*s", mark, status, pathWidth, path)

			if i == state.ActiveFile {
				if state.MarkedFiles[file.Path] {
					line = MarkedStyle.Width(w - 4).Render(line)
				} else {
					line = SelectedStyle.Width(w - 4).Render(line)
				}
			}

			lines = append(lines, line)
		}
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)

	style := PanelBorder.Width(w).Height(h)
	if state.Focus == enums.PanelFiles {
		style = FocusedBorder.Width(w).Height(h)
	}

	return style.Render(content)
}
