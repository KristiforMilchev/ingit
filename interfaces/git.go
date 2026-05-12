package interfaces

import "ingit/models"

type GitProvider interface {
	DiscoverRepos() []models.Repo
	ReadGitState(repo models.Repo) models.GitState
	ReadDiffForFile(repo models.Repo, file models.FileStatus) string
	ReadBranchGraph(repo models.Repo) string
	ReadRecentCommits(repo models.Repo) []models.GitCommit
	StageFile(repo models.Repo, path string) error
	UnstageFile(repo models.Repo, path string) error
	StageAll(repo models.Repo) error
	UnstageAll(repo models.Repo) error
	ExecutePushPlan(repo models.Repo, plan models.PushPlanExecution) error
	ExecuteQuickPush(repo models.Repo, message string) error
	GitOutput(path string, args ...string) (string, error)
	GitOutputWithInput(repoPath string, input string, args ...string) (string, error)
}
