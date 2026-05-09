package views

import (
	"strings"

	"ingit/utils"
)

type DiffRow struct {
	Left      string
	Right     string
	LeftKind  string
	RightKind string
}

func IsRealDiff(value string) bool {
	return strings.HasPrefix(value, "diff --git") || strings.Contains(value, "\n@@")
}

func RenderSideBySideDiff(diff string, columnWidth int) string {
	var rows []DiffRow
	lines := strings.Split(diff, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") ||
			strings.HasPrefix(line, "index ") ||
			strings.HasPrefix(line, "---") ||
			strings.HasPrefix(line, "+++") {
			continue
		}

		if strings.HasPrefix(line, "@@") {
			clean := strings.TrimSpace(line)
			rows = append(rows, DiffRow{Left: clean, Right: clean, LeftKind: "hunk", RightKind: "hunk"})
			continue
		}

		if strings.HasPrefix(line, "-") {
			rows = append(rows, DiffRow{Left: line, Right: "", LeftKind: "removed", RightKind: "empty"})
			continue
		}

		if strings.HasPrefix(line, "+") {
			rows = append(rows, DiffRow{Left: "", Right: line, LeftKind: "empty", RightKind: "added"})
			continue
		}

		if strings.TrimSpace(line) == "" {
			rows = append(rows, DiffRow{Left: "", Right: "", LeftKind: "context", RightKind: "context"})
			continue
		}

		rows = append(rows, DiffRow{Left: line, Right: line, LeftKind: "context", RightKind: "context"})
	}

	out := []string{
		TitleStyle.Render(utils.PadRight("OLD", columnWidth)) + " │ " + TitleStyle.Render("NEW"),
		MutedStyle.Render(strings.Repeat("─", columnWidth)) + "─┼─" + MutedStyle.Render(strings.Repeat("─", columnWidth)),
	}

	for _, row := range rows {
		left := RenderDiffCell(row.Left, row.LeftKind, columnWidth)
		right := RenderDiffCell(row.Right, row.RightKind, columnWidth)
		out = append(out, left+" │ "+right)
	}

	return strings.Join(out, "\n")
}

func RenderDiffCell(value string, kind string, width int) string {
	value = utils.TruncateLine(value, width)
	padded := utils.PadRight(value, width)

	switch kind {
	case "removed":
		return ErrorStyle.Render(padded)
	case "added":
		return OkStyle.Render(padded)
	case "hunk":
		return TitleStyle.Render(padded)
	case "empty":
		return MutedStyle.Render(padded)
	default:
		return MutedStyle.Render(padded)
	}
}

func RenderUnifiedDiff(diff string, width int) string {
	var out []string
	for _, line := range strings.Split(diff, "\n") {
		line = utils.TruncateLine(line, width)
		switch {
		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			out = append(out, OkStyle.Render(line))
		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---"):
			out = append(out, ErrorStyle.Render(line))
		case strings.HasPrefix(line, "@@"):
			out = append(out, TitleStyle.Render(line))
		default:
			out = append(out, MutedStyle.Render(line))
		}
	}
	return strings.Join(out, "\n")
}
