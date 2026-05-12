package models

type DiffHunk struct {
	FilePath string
	Header   string
	Lines    []string
	Selected bool
}
