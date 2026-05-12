package enums

type Mode int

const (
	ModeGit Mode = iota
	ModePush
	ModeQuickPush
	ModeBranches
	ModeSurgery
)
