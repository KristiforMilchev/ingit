package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"ingit/implementations"
	"ingit/views"
)

func main() {
	git := implementations.NewGitCLI()
	repos := git.DiscoverRepos()

	if len(repos) == 0 {
		fmt.Println("No git repositories found. Run this from a directory containing git repos.")
		os.Exit(1)
	}

	model := views.NewModel(git, repos)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
