package interfaces

type PlannedBranch struct {
	Name    string
	Base    string
	Message string
	Status  string
	Error   string
	Files   []string
	Commits []PlannedCommit
}

type PlannedCommit struct {
	Message string
	Files   []string
}

type ChangePlan struct {
	Branches []PlannedBranch
}
