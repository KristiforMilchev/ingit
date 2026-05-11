package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"ingit/states"
	"ingit/utils"
)

func RenderBranchesView(state states.BranchesState, width int, height int) string {
	if state.Error != "" {
		content := TitleStyle.Render("Branches") +
			"\n\n" +
			ErrorStyle.Render(state.Error) +
			"\n\n" +
			MutedStyle.Render("enter checkout · m merge · M merge+delete · d delete · tab details · esc back")

		return PanelBorder.
			Width(width - 2).
			Height(height).
			Render(content)
	}

	leftW := 44
	if width >= 170 {
		leftW = 56
	}

	rightW := width - leftW - 3
	if rightW < 60 {
		rightW = 60
	}

	left := renderBranchList(state, leftW, height)
	right := renderBranchDetails(state, rightW, height)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func renderBranchList(state states.BranchesState, w int, h int) string {
	lines := []string{
		TitleStyle.Render("Branches"),
		MutedStyle.Render("enter checkout · m merge · M merge+delete · d delete · tab details · esc back"),
		"",
	}

	if len(state.Branches) == 0 {
		lines = append(lines, MutedStyle.Render("no branches"))
	} else {
		for i, branch := range state.Branches {
			current := " "
			if branch.Current {
				current = "*"
			}

			remote := ""
			if branch.Remote {
				remote = " remote"
			}

			line := fmt.Sprintf("%s %s%s", current, branch.Name, remote)
			line = utils.TruncateLine(line, w-4)

			if i == state.Selected {
				line = SelectedStyle.Width(w - 4).Render(line)
			} else if branch.Current {
				line = OkStyle.Render(line)
			} else if branch.Remote {
				line = WarnStyle.Render(line)
			} else {
				line = MutedStyle.Render(line)
			}

			lines = append(lines, line)
		}
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)

	return FocusedBorder.
		Width(w).
		Height(h).
		Render(content)
}

func renderBranchDetails(state states.BranchesState, w int, h int) string {
	branch := selectedBranchName(state)
	if branch == "" {
		return PanelBorder.
			Width(w).
			Height(h).
			Render(MutedStyle.Render("no branch selected"))
	}

	title := "Commits"
	if state.Panel == 1 {
		title = "Files"
	}

	lines := []string{
		TitleStyle.Render(title) + "  " + MutedStyle.Render(branch),
		"",
	}

	if state.Panel == 0 {
		if len(state.Commits) == 0 {
			lines = append(lines, MutedStyle.Render("no commits"))
		} else {
			for _, commit := range state.Commits {
				lines = append(lines, renderGraphLine(utils.TruncateLine(commit, w-4)))
			}
		}
	} else {
		if len(state.Files) == 0 {
			lines = append(lines, MutedStyle.Render("no files"))
		} else {
			for _, file := range state.Files {
				lines = append(lines, MutedStyle.Render(utils.TruncateLine(file, w-4)))
			}
		}
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)

	return PanelBorder.
		Width(w).
		Height(h).
		Render(content)
}

func selectedBranchName(state states.BranchesState) string {
	if len(state.Branches) == 0 ||
		state.Selected < 0 ||
		state.Selected >= len(state.Branches) {
		return ""
	}

	return state.Branches[state.Selected].Name
}
