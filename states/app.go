package states

import (
	"ingit/enums"
	"ingit/interfaces"
	"ingit/models"
)

type AppState struct {
	Repos        []models.Repo
	States       map[string]models.GitState
	ActiveRepo   int
	ActiveFile   int
	Focus        enums.Panel
	Mode         enums.Mode
	GraphMode    enums.GraphMode
	Plan         interfaces.ChangePlan
	ActiveBranch int
	MarkedFiles  map[string]bool

	NamingBranch      bool
	NamingMessage     bool
	QuickPush         bool
	BranchInput       string
	MessageInput      string
	QuickMessageInput string

	Width     int
	Height    int
	StatusMsg string
}
