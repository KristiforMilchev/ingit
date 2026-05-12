package states

import "ingit/models"

type BranchesState struct {
	RepoPath string

	Branches []models.GitBranch
	Selected int

	CurrentBranch string
	Panel         int

	Commits []string
	Files   []string

	Error string

	Merged     bool
	Deleted    bool
	Done       bool
	CheckedOut bool

	Creating    bool
	BranchInput string
}
