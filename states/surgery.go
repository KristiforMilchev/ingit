package states

import "ingit/models"

type SurgeryState struct {
	RepoPath string
	FilePath string

	Hunks    []models.DiffHunk
	Selected int

	MessageInput   string
	WritingMessage bool

	Error string
	Done  bool

	Executed bool
}
