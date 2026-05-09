package models

type GitState struct {
	Branch  string
	Ahead   int
	Behind  int
	Files   []FileStatus
	Diff    string
	Graph   string
	Commits []GitCommit
	Error   string
}
