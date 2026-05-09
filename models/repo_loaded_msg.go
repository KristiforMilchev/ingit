package models

type RepoLoadedMsg struct {
	Index int
	State GitState
}
