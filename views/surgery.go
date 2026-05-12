package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"ingit/states"
	"ingit/utils"
)

func RenderSurgeryView(state states.SurgeryState, width int, height int) string {
	if state.WritingMessage {
		return renderSurgeryMessagePopup(state, width, height)
	}

	leftW := 48
	if width >= 160 {
		leftW = 62
	}

	rightW := width - leftW - 3
	if rightW < 60 {
		rightW = 60
	}

	left := renderSurgeryHunks(state, leftW, height)
	right := renderSurgeryPreview(state, rightW, height)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func renderSurgeryHunks(state states.SurgeryState, w int, h int) string {
	lines := []string{
		TitleStyle.Render("Commit Surgery"),
		MutedStyle.Render("space select · C message · e execute · esc back"),
		MutedStyle.Render(state.FilePath),
		"",
	}

	if state.Error != "" {
		lines = append(lines, ErrorStyle.Render(utils.TruncateLine(state.Error, w-4)), "")
	}

	if len(state.Hunks) == 0 {
		lines = append(lines, MutedStyle.Render("no hunks"))
	} else {
		for i, hunk := range state.Hunks {
			mark := "[ ]"
			if hunk.Selected {
				mark = "[x]"
			}

			line := fmt.Sprintf("%s %s", mark, hunk.Header)
			line = utils.TruncateLine(line, w-4)

			if i == state.Selected {
				line = SelectedStyle.Width(w - 4).Render(line)
			} else if hunk.Selected {
				line = OkStyle.Render(line)
			} else {
				line = MutedStyle.Render(line)
			}

			lines = append(lines, line)
		}
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)
	return FocusedBorder.Width(w).Height(h).Render(content)
}

func renderSurgeryPreview(state states.SurgeryState, w int, h int) string {
	lines := []string{
		TitleStyle.Render("Selected Patch"),
		"",
	}

	for _, hunk := range state.Hunks {
		if !hunk.Selected {
			continue
		}

		lines = append(lines, TitleStyle.Render(utils.TruncateLine(hunk.Header, w-4)))

		for _, line := range hunk.Lines {
			lines = append(lines, renderGraphLine(utils.TruncateLine(line, w-4)))
		}

		lines = append(lines, "")
	}

	if len(lines) == 2 {
		lines = append(lines, MutedStyle.Render("no hunks selected"))
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)
	return PanelBorder.Width(w).Height(h).Render(content)
}

func renderSurgeryMessagePopup(state states.SurgeryState, width int, height int) string {
	popupW := 64
	if width < 80 {
		popupW = width - 8
	}
	if popupW < 36 {
		popupW = 36
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		TitleStyle.Render("Commit Selected Hunks"),
		"",
		MutedStyle.Render("Commit message"),
		InputStyle.Render(utils.TruncateLine(" "+state.MessageInput+" ", popupW-6)),
		"",
		MutedStyle.Render("enter commit · esc cancel"),
	)

	popup := FocusedBorder.Width(popupW).Render(content)

	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		popup,
	)
}
