package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"ingit/enums"
	"ingit/states"
	"ingit/utils"
)

func renderCurrentCommitLines(state states.AppState, width int) []string {
	repoState := currentState(state)

	current := repoState.Branch
	if current == "" {
		current = "unknown"
	}

	lines := []string{
		MutedStyle.Render(utils.TruncateLine("base: "+current, width)),
	}

	if len(repoState.Commits) > 0 {
		commit := repoState.Commits[0]
		line := "head: " + commit.Hash

		if commit.Refs != "" {
			line += " " + commit.Refs
		}

		if commit.Message != "" {
			line += " " + commit.Message
		}

		lines = append(lines, MutedStyle.Render(utils.TruncateLine(line, width)))
	}

	return lines
}

func renderPushPlans(state states.AppState, w, h int) string {
	lines := []string{TitleStyle.Render("Push Plans")}
	lines = append(lines, renderCurrentCommitLines(state, w-4)...)
	lines = append(lines, "")

	if state.NamingBranch {
		label := "new"
		if state.QuickPush {
			label = "quick"
		}

		line := fmt.Sprintf("%s branch: %s", label, state.BranchInput)
		lines = append(lines, InputStyle.Render(utils.TruncateLine(" "+line+" ", w-4)))
		lines = append(lines, MutedStyle.Render("enter create · esc cancel"))
	}

	if state.NamingMessage {
		line := "message: " + state.MessageInput
		lines = append(lines, InputStyle.Render(utils.TruncateLine(" "+line+" ", w-4)))
		lines = append(lines, MutedStyle.Render("enter set · esc cancel"))
	}

	if len(state.Plan.Branches) == 0 {
		lines = append(lines, MutedStyle.Render("no push plans"))
		lines = append(lines, MutedStyle.Render("press b to create one"))
	} else {
		for i, branch := range state.Plan.Branches {
			style := planStyle(i)
			line := fmt.Sprintf("■ %s", branch.Name)

			if i == state.ActiveBranch {
				line = SelectedStyle.Width(w - 4).Render(line)
			} else {
				line = style.Render(utils.TruncateLine(line, w-4))
			}

			lines = append(lines, line)

			base := fmt.Sprintf("  base: %s", branch.Base)
			lines = append(lines, MutedStyle.Render(utils.TruncateLine(base, w-6)))

			status := branch.Status
			if status == "" {
				status = "pending"
			}

			lines = append(lines, MutedStyle.Render(utils.TruncateLine("  status: "+status, w-6)))

			message := branch.Message
			if strings.TrimSpace(message) == "" {
				message = "<missing commit message>"
				lines = append(lines, WarnStyle.Render(utils.TruncateLine("  msg: "+message, w-6)))
			} else {
				lines = append(lines, OkStyle.Render(utils.TruncateLine("  msg: "+message, w-6)))
			}

			if branch.Error != "" {
				lines = append(lines, ErrorStyle.Render(utils.TruncateLine("  error: "+branch.Error, w-6)))
			}

			if len(branch.Files) == 0 {
				lines = append(lines, MutedStyle.Render("  no files assigned"))
			} else {
				for _, file := range branch.Files {
					lines = append(lines, style.Render(utils.TruncateLine("  • "+file, w-6)))
				}
			}
		}
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)

	style := PanelBorder.Width(w).Height(h)
	if state.Focus == enums.PanelPlans {
		style = FocusedBorder.Width(w).Height(h)
	}

	return style.Render(content)
}

func renderAssignmentFiles(state states.AppState, w, h int) string {
	repoState := currentState(state)

	lines := []string{
		TitleStyle.Render("Changed Files") + "  " + MutedStyle.Render("assigned files stay visible"),
	}

	if repoState.Error != "" {
		lines = append(lines, ErrorStyle.Render(repoState.Error))
	} else if len(repoState.Files) == 0 {
		lines = append(lines, MutedStyle.Render("clean"))
	} else {
		pathWidth := w - 10
		if pathWidth < 14 {
			pathWidth = 14
		}

		for i, file := range repoState.Files {
			status := fmt.Sprintf("%s%s", normalizeStatus(file.Index), normalizeStatus(file.Worktree))

			mark := " "
			if state.MarkedFiles[file.Path] {
				mark = "●"
			}

			branchIndex, branchName, assigned := assignedBranchForFile(state, file.Path)

			label := file.Path
			if assigned {
				label = fmt.Sprintf("%s → %s", file.Path, branchName)
			}

			line := fmt.Sprintf("%s %-2s %-*s", mark, status, pathWidth, utils.TruncateLine(label, pathWidth))

			if assigned {
				line = planStyle(branchIndex).Render(utils.TruncateLine(line, w-4))
			}

			if i == state.ActiveFile {
				if state.MarkedFiles[file.Path] {
					line = MarkedStyle.Width(w - 4).Render(line)
				} else {
					line = SelectedStyle.Width(w - 4).Render(line)
				}
			}

			lines = append(lines, line)
		}

		lines = append(lines, "")
		lines = append(lines, TitleStyle.Render("Legend"))

		for i, branch := range state.Plan.Branches {
			lines = append(lines, planStyle(i).Render(utils.TruncateLine("■ "+branch.Name, w-4)))
		}
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)

	style := PanelBorder.Width(w).Height(h)
	if state.Focus == enums.PanelFiles {
		style = FocusedBorder.Width(w).Height(h)
	}

	return style.Render(content)
}

func assignedBranchForFile(state states.AppState, path string) (int, string, bool) {
	for i, branch := range state.Plan.Branches {
		for _, file := range branch.Files {
			if file == path {
				return i, branch.Name, true
			}
		}
	}

	return -1, "", false
}

func planStyle(index int) lipgloss.Style {
	switch index % 6 {
	case 0:
		return PlanColor1
	case 1:
		return PlanColor2
	case 2:
		return PlanColor3
	case 3:
		return PlanColor4
	case 4:
		return PlanColor5
	default:
		return PlanColor6
	}
}
