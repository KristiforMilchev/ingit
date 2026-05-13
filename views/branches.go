package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"ingit/states"
	"ingit/utils"
)

type BranchesView struct {
	state  states.BranchesState
	width  int
	height int
	render string
}

func RenderBranchesView(state states.BranchesState, width int, height int) string {
	view := BranchesView{
		state:  state,
		width:  width,
		height: height,
	}
	if view.state.ShowPopup {
		return view.popupView()
	}

	view.renderBranchesBase()

	if state.Creating {
		view.createNewBranch()
	}

	return view.render
}

func (b *BranchesView) renderBranchesBase() {
	if b.state.Error != "" {
		content := TitleStyle.Render("Branches") +
			"\n\n" +
			ErrorStyle.Render(b.state.Error) +
			"\n\n" +
			MutedStyle.Render("enter checkout · m merge · M merge+delete · d delete · tab details · p - push current · esc back")

		b.render = PanelBorder.
			Width(b.width - 2).
			Height(b.height).
			Render(content)
		return
	}

	leftW := 44
	if b.width >= 170 {
		leftW = 56
	}

	rightW := b.width - leftW - 3
	if rightW < 60 {
		rightW = 60
	}

	left := renderBranchList(b.state, leftW, b.height)
	right := renderBranchDetails(b.state, rightW, b.height)

	b.render = lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (b *BranchesView) createNewBranch() {
	popupW := 56
	if b.width < 70 {
		popupW = b.width - 8
	}
	if popupW < 32 {
		popupW = 32
	}

	input := InputStyle.Render(
		utils.TruncateLine(" "+b.state.BranchInput+" ", popupW-6),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		TitleStyle.Render("Create Branch"),
		"",
		MutedStyle.Render("Branch name"),
		input,
		"",
		MutedStyle.Render("enter create · esc cancel"),
	)

	popup := FocusedBorder.
		Width(popupW).
		Render(content)

	b.render = lipgloss.Place(
		b.width,
		b.height,
		lipgloss.Center,
		lipgloss.Center,
		popup,
		lipgloss.WithWhitespaceChars(" "),
	)
}

func (b *BranchesView) popupView() string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Width(80)

	content := fmt.Sprintf(
		"%s\n\n%s\n\n[esc] close",
		b.state.PopupTitle,
		b.state.PopupBody,
	)

	return lipgloss.Place(
		b.width,
		b.height,
		lipgloss.Center,
		lipgloss.Center,
		box.Render(content),
	)
}

func renderBranchList(state states.BranchesState, w int, h int) string {
	lines := []string{
		TitleStyle.Render("Branches"),
		MutedStyle.Render("enter checkout · m merge · M merge+delete · d delete · tab details · n new branch · esc back"),
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
