package models

type CommitSlice struct {
	Message string
	File    string
	Hunks   []DiffHunk
}
