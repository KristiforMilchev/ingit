package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"ingit/controllers"
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

	model := controllers.NewModel(git, repos, views.RenderAppView)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
