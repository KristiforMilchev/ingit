package views

import (
	"ingit/models"
	"ingit/states"
)

func currentState(state states.AppState) models.GitState {
	if len(state.Repos) == 0 {
		return models.GitState{}
	}

	if state.ActiveRepo < 0 || state.ActiveRepo >= len(state.Repos) {
		return models.GitState{}
	}

	return state.States[state.Repos[state.ActiveRepo].Path]
}

func selectedFile(state states.AppState) (models.FileStatus, bool) {
	files := currentState(state).Files

	if len(files) == 0 || state.ActiveFile < 0 || state.ActiveFile >= len(files) {
		return models.FileStatus{}, false
	}

	return files[state.ActiveFile], true
}

func selectedDiffModeLabel(state states.AppState) string {
	file, ok := selectedFile(state)
	if !ok {
		return "diff"
	}

	if hasWorktreeChange(file) {
		return "unstaged diff"
	}

	if hasStagedChange(file) {
		return "staged diff"
	}

	return "diff"
}

func hasStagedChange(file models.FileStatus) bool {
	return file.Index != "" && file.Index != " " && file.Index != "." && file.Index != "?"
}

func hasWorktreeChange(file models.FileStatus) bool {
	return file.Worktree != "" && file.Worktree != " " && file.Worktree != "."
}

func normalizeStatus(s string) string {
	if s == " " || s == "" {
		return "."
	}

	return s
}
