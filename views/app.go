package views

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"ingit/enums"
	"ingit/interfaces"
	"ingit/models"
	"ingit/utils"
)

type Model struct {
	Git interfaces.GitProvider

	Repos        []models.Repo
	States       map[string]models.GitState
	ActiveRepo   int
	ActiveFile   int
	Focus        enums.Panel
	Mode         enums.Mode
	GraphMode    enums.GraphMode
	Branches     BranchesModel
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

func NewModel(git interfaces.GitProvider, repos []models.Repo) Model {
	return Model{
		Git:         git,
		Repos:       repos,
		States:      map[string]models.GitState{},
		ActiveRepo:  0,
		ActiveFile:  0,
		Focus:       enums.PanelFiles,
		Mode:        enums.ModeGit,
		GraphMode:   enums.GraphVisual,
		StatusMsg:   "ready",
		MarkedFiles: map[string]bool{},
		Plan:        interfaces.ChangePlan{Branches: []interfaces.PlannedBranch{}},
		Branches:    NewBranchesModel(),
	}
}

func (m Model) Init() tea.Cmd {
	if len(m.Repos) == 0 {
		return tickCmd()
	}

	return tea.Batch(loadRepoCmd(m.Git, 0, m.Repos[0]), tickCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case models.TickMsg:
		if len(m.Repos) == 0 {
			return m, tickCmd()
		}
		return m, tea.Batch(loadRepoCmd(m.Git, m.ActiveRepo, m.Repos[m.ActiveRepo]), tickCmd())

	case models.RepoLoadedMsg:
		if msg.Index < 0 || msg.Index >= len(m.Repos) {
			return m, nil
		}

		repo := m.Repos[msg.Index]
		m.States[repo.Path] = msg.State

		if m.StatusMsg == "" || strings.HasPrefix(m.StatusMsg, "refreshed") {
			m.StatusMsg = "refreshed " + repo.Name
		}

		return m, nil

	case models.ActionDoneMsg:
		if len(m.Repos) == 0 {
			return m, nil
		}

		if msg.Err != nil {
			m.StatusMsg = msg.Err.Error()
		} else {
			m.StatusMsg = msg.Message
		}

		return m, loadRepoCmd(m.Git, m.ActiveRepo, m.Repos[m.ActiveRepo])

	case models.PlanExecutedMsg:
		if len(m.Repos) == 0 {
			return m, nil
		}

		if msg.Index >= 0 && msg.Index < len(m.Plan.Branches) {
			if msg.Err != nil {
				m.Plan.Branches[msg.Index].Status = "failed"
				m.Plan.Branches[msg.Index].Error = msg.Err.Error()
				m.StatusMsg = msg.Err.Error()
			} else {
				executed := m.Plan.Branches[msg.Index].Name
				m.Plan.Branches = append(m.Plan.Branches[:msg.Index], m.Plan.Branches[msg.Index+1:]...)
				m.MarkedFiles = map[string]bool{}

				if len(m.Plan.Branches) == 0 {
					m.ActiveBranch = 0
				} else if m.ActiveBranch >= len(m.Plan.Branches) {
					m.ActiveBranch = len(m.Plan.Branches) - 1
				}

				m.StatusMsg = "executed and cleared " + executed
			}
		}

		return m, loadRepoCmd(m.Git, m.ActiveRepo, m.Repos[m.ActiveRepo])

	case models.QuickPushExecutedMsg:
		if len(m.Repos) == 0 {
			return m, nil
		}

		if msg.Err != nil {
			m.StatusMsg = msg.Err.Error()
			return m, loadRepoCmd(m.Git, m.ActiveRepo, m.Repos[m.ActiveRepo])
		}

		m.Mode = enums.ModeGit
		m.Focus = enums.PanelFiles
		m.QuickMessageInput = ""
		m.MarkedFiles = map[string]bool{}
		m.StatusMsg = "quick push completed"

		return m, loadRepoCmd(m.Git, m.ActiveRepo, m.Repos[m.ActiveRepo])

	case models.AllPlansExecutedMsg:
		if len(m.Repos) == 0 {
			return m, nil
		}

		failed := 0
		failedByIndex := map[int]error{}

		for _, result := range msg.Results {
			if result.Err != nil {
				failed++
				failedByIndex[result.Index] = result.Err
			}
		}

		var remaining []interfaces.PlannedBranch
		for index, branch := range m.Plan.Branches {
			if err, ok := failedByIndex[index]; ok {
				branch.Status = "failed"
				branch.Error = err.Error()
				remaining = append(remaining, branch)
			}
		}

		m.Plan.Branches = remaining
		m.MarkedFiles = map[string]bool{}

		if len(m.Plan.Branches) == 0 {
			m.ActiveBranch = 0
		} else if m.ActiveBranch >= len(m.Plan.Branches) {
			m.ActiveBranch = len(m.Plan.Branches) - 1
		}

		if failed > 0 {
			m.StatusMsg = fmt.Sprintf("executed all; %d failed plan(s) kept", failed)
		} else {
			m.StatusMsg = "all push plans executed and cleared"
		}

		return m, loadRepoCmd(m.Git, m.ActiveRepo, m.Repos[m.ActiveRepo])

	case tea.KeyMsg:
		if len(m.Repos) == 0 {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			}
			return m, nil
		}

		if m.Mode == enums.ModeBranches {
			var cmd tea.Cmd
			m.Branches, cmd = m.Branches.Update(msg)

			if m.Branches.Done {
				m.Branches.Done = false
				m.Mode = enums.ModeGit
				m.Focus = enums.PanelFiles
				m.StatusMsg = "git view"
				return m, nil
			}

			if m.Branches.CheckedOut {
				branch := m.Branches.SelectedBranch()
				m.Branches.CheckedOut = false
				m.Mode = enums.ModeGit
				m.Focus = enums.PanelFiles

				if branch == "" {
					m.StatusMsg = "branch switched"
				} else {
					m.StatusMsg = "switched to " + branch
				}

				return m, loadRepoCmd(m.Git, m.ActiveRepo, m.Repos[m.ActiveRepo])
			}

			return m, cmd
		}

		if m.Mode == enums.ModeQuickPush {
			return m.handleQuickPushInput(msg)
		}

		if m.NamingBranch {
			return m.handleBranchNameInput(msg)
		}

		if m.NamingMessage {
			return m.handleMessageInput(msg)
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "p":
			m.Mode = enums.ModePush
			m.Focus = enums.PanelPlans
			m.StatusMsg = "push planner"
			return m, nil

		case "P":
			m.startQuickPushView()
			return m, nil

		case "B":
			repo := m.Repos[m.ActiveRepo]
			m.Branches.Load(repo.Path)
			m.Mode = enums.ModeBranches
			m.Focus = enums.PanelGraph
			m.StatusMsg = "branches"
			return m, nil

		case "esc":
			m.Mode = enums.ModeGit
			m.Focus = enums.PanelFiles
			m.StatusMsg = "git view"
			return m, nil

		case "g":
			if m.GraphMode == enums.GraphVisual {
				m.GraphMode = enums.GraphText
				m.StatusMsg = "graph: branch text"
			} else {
				m.GraphMode = enums.GraphVisual
				m.StatusMsg = "graph: git log"
			}
			return m, nil

		case "tab", "l", "right":
			m.nextPanel()
			return m, nil

		case "h", "left":
			m.prevPanel()
			return m, nil

		case "j", "down":
			m.moveDown()
			return m, m.loadDiffCmd()

		case "k", "up":
			m.moveUp()
			return m, m.loadDiffCmd()

		case "r":
			m.StatusMsg = "refreshing..."
			return m, loadRepoCmd(m.Git, m.ActiveRepo, m.Repos[m.ActiveRepo])

		case " ":
			m.toggleMarkedFile()
			return m, nil

		case "b":
			m.startBranchNameInput()
			return m, nil

		case "C":
			m.startMessageInput()
			return m, nil

		case "a":
			m.assignMarkedToBranch()
			return m, nil

		case "x":
			m.removeSelectedFileFromPlan()
			return m, nil

		case "e":
			return m, m.executeSelectedPlanCmd()

		case "E":
			return m, m.executeAllPlansCmd()

		case "s":
			if m.Mode == enums.ModePush {
				m.StatusMsg = "staging is disabled in push planner"
				return m, nil
			}
			return m, m.stageSelectedCmd()

		case "u":
			if m.Mode == enums.ModePush {
				m.StatusMsg = "unstaging is disabled in push planner"
				return m, nil
			}
			return m, m.unstageSelectedCmd()

		case "S":
			if m.Mode == enums.ModePush {
				m.StatusMsg = "stage all is disabled in push planner"
				return m, nil
			}
			return m, m.stageAllCmd()

		case "U":
			if m.Mode == enums.ModePush {
				m.StatusMsg = "unstage all is disabled in push planner"
				return m, nil
			}
			return m, m.unstageAllCmd()
		}
	}

	return m, nil
}

func (m Model) handleBranchNameInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.NamingBranch = false
		m.QuickPush = false
		m.BranchInput = ""
		m.StatusMsg = "branch creation cancelled"
		return m, nil

	case "enter":
		name := strings.TrimSpace(m.BranchInput)
		if name == "" {
			m.StatusMsg = "branch name cannot be empty"
			return m, nil
		}

		branch := interfaces.PlannedBranch{
			Name:   name,
			Base:   m.currentState().Branch,
			Status: "pending",
		}

		if m.QuickPush {
			for _, file := range m.currentState().Files {
				branch.Files = append(branch.Files, file.Path)
				m.removeFileFromAllBranches(file.Path)
			}
		}

		m.Plan.Branches = append(m.Plan.Branches, branch)
		m.ActiveBranch = len(m.Plan.Branches) - 1
		m.NamingBranch = false
		m.QuickPush = false
		m.BranchInput = ""

		if len(branch.Files) > 0 {
			m.StatusMsg = fmt.Sprintf("quick push plan %s with %d file(s)", name, len(branch.Files))
		} else {
			m.StatusMsg = "created push plan " + name
		}

		return m, nil

	case "backspace":
		if len(m.BranchInput) > 0 {
			runes := []rune(m.BranchInput)
			m.BranchInput = string(runes[:len(runes)-1])
		}
		return m, nil

	default:
		if len(msg.Runes) > 0 {
			m.BranchInput += string(msg.Runes)
		}
		return m, nil
	}
}

func (m Model) handleMessageInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.NamingMessage = false
		m.MessageInput = ""
		m.StatusMsg = "commit message cancelled"
		return m, nil

	case "enter":
		if len(m.Plan.Branches) == 0 {
			m.NamingMessage = false
			m.MessageInput = ""
			m.StatusMsg = "no push plan selected"
			return m, nil
		}

		message := strings.TrimSpace(m.MessageInput)
		if message == "" {
			m.StatusMsg = "commit message cannot be empty"
			return m, nil
		}

		if m.ActiveBranch < 0 || m.ActiveBranch >= len(m.Plan.Branches) {
			m.ActiveBranch = 0
		}

		m.Plan.Branches[m.ActiveBranch].Message = message
		m.NamingMessage = false
		m.MessageInput = ""
		m.StatusMsg = "message set for " + m.Plan.Branches[m.ActiveBranch].Name

		return m, nil

	case "backspace":
		if len(m.MessageInput) > 0 {
			runes := []rune(m.MessageInput)
			m.MessageInput = string(runes[:len(runes)-1])
		}
		return m, nil

	default:
		if len(msg.Runes) > 0 {
			m.MessageInput += string(msg.Runes)
		}
		return m, nil
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(900*time.Millisecond, func(t time.Time) tea.Msg {
		return models.TickMsg(t)
	})
}

func (m *Model) nextPanel() {
	if m.Mode == enums.ModePush {
		switch m.Focus {
		case enums.PanelPlans:
			m.Focus = enums.PanelGraph
		case enums.PanelGraph:
			m.Focus = enums.PanelFiles
		default:
			m.Focus = enums.PanelPlans
		}
		return
	}

	if len(m.Repos) == 1 {
		m.Focus = enums.PanelFiles
		return
	}

	if m.Focus == enums.PanelRepos {
		m.Focus = enums.PanelFiles
	} else {
		m.Focus = enums.PanelRepos
	}
}

func (m *Model) prevPanel() {
	if m.Mode == enums.ModePush {
		switch m.Focus {
		case enums.PanelPlans:
			m.Focus = enums.PanelFiles
		case enums.PanelGraph:
			m.Focus = enums.PanelPlans
		default:
			m.Focus = enums.PanelGraph
		}
		return
	}

	m.nextPanel()
}

func (m *Model) moveDown() {
	switch m.Focus {
	case enums.PanelRepos:
		if m.ActiveRepo < len(m.Repos)-1 {
			m.ActiveRepo++
			m.ActiveFile = 0
		}
	case enums.PanelFiles:
		files := m.currentState().Files
		if m.ActiveFile < len(files)-1 {
			m.ActiveFile++
		}
	case enums.PanelPlans:
		if m.ActiveBranch < len(m.Plan.Branches)-1 {
			m.ActiveBranch++
		}
	}
}

func (m *Model) moveUp() {
	switch m.Focus {
	case enums.PanelRepos:
		if m.ActiveRepo > 0 {
			m.ActiveRepo--
			m.ActiveFile = 0
		}
	case enums.PanelFiles:
		if m.ActiveFile > 0 {
			m.ActiveFile--
		}
	case enums.PanelPlans:
		if m.ActiveBranch > 0 {
			m.ActiveBranch--
		}
	}
}

func (m *Model) startQuickPushView() {
	m.Mode = enums.ModeQuickPush
	m.Focus = enums.PanelFiles
	m.QuickMessageInput = ""
	m.StatusMsg = "quick push: type commit message, enter to commit+push, esc to cancel"
}

func (m Model) handleQuickPushInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.Mode = enums.ModeGit
		m.Focus = enums.PanelFiles
		m.QuickMessageInput = ""
		m.StatusMsg = "quick push cancelled"
		return m, nil

	case "enter":
		message := strings.TrimSpace(m.QuickMessageInput)
		if message == "" {
			m.StatusMsg = "commit message cannot be empty"
			return m, nil
		}

		return m, m.executeQuickPushCmd(message)

	case "backspace":
		if len(m.QuickMessageInput) > 0 {
			runes := []rune(m.QuickMessageInput)
			m.QuickMessageInput = string(runes[:len(runes)-1])
		}
		return m, nil

	default:
		if len(msg.Runes) > 0 {
			m.QuickMessageInput += string(msg.Runes)
		}
		return m, nil
	}
}

func (m *Model) startBranchNameInput() {
	m.Mode = enums.ModePush
	m.Focus = enums.PanelPlans
	m.NamingBranch = true
	m.QuickPush = false
	m.BranchInput = "feature/"
	m.StatusMsg = "enter push plan branch name"
}

func (m *Model) startQuickPushInput() {
	m.Mode = enums.ModePush
	m.Focus = enums.PanelPlans
	m.NamingBranch = true
	m.QuickPush = true
	m.BranchInput = "feature/"
	m.StatusMsg = "quick push: name branch, all changed files will be assigned"
}

func (m *Model) startMessageInput() {
	m.Mode = enums.ModePush
	m.Focus = enums.PanelPlans

	if len(m.Plan.Branches) == 0 {
		m.StatusMsg = "create a push plan first"
		return
	}

	if m.ActiveBranch < 0 || m.ActiveBranch >= len(m.Plan.Branches) {
		m.ActiveBranch = 0
	}

	m.NamingMessage = true
	m.MessageInput = m.Plan.Branches[m.ActiveBranch].Message
	m.StatusMsg = "enter commit message for " + m.Plan.Branches[m.ActiveBranch].Name
}

func (m *Model) toggleMarkedFile() {
	file, ok := m.selectedFile()
	if !ok {
		return
	}

	key := file.Path
	if m.MarkedFiles[key] {
		delete(m.MarkedFiles, key)
		m.StatusMsg = "unmarked " + key
	} else {
		m.MarkedFiles[key] = true
		m.StatusMsg = "marked " + key
	}
}

func (m *Model) assignMarkedToBranch() {
	if len(m.Plan.Branches) == 0 {
		m.StatusMsg = "create a push plan first with b"
		return
	}

	if m.ActiveBranch < 0 || m.ActiveBranch >= len(m.Plan.Branches) {
		m.ActiveBranch = 0
	}

	branch := &m.Plan.Branches[m.ActiveBranch]
	count := 0

	for path := range m.MarkedFiles {
		m.removeFileFromAllBranches(path)

		if !contains(branch.Files, path) {
			branch.Files = append(branch.Files, path)
			count++
		}

		delete(m.MarkedFiles, path)
	}

	m.StatusMsg = fmt.Sprintf("assigned %d file(s) to %s", count, branch.Name)
}

func (m *Model) removeSelectedFileFromPlan() {
	file, ok := m.selectedFile()
	if !ok {
		return
	}

	m.removeFileFromAllBranches(file.Path)
	delete(m.MarkedFiles, file.Path)
	m.StatusMsg = "removed from push plan: " + file.Path
}

func (m *Model) removeFileFromAllBranches(path string) {
	for bi := range m.Plan.Branches {
		files := m.Plan.Branches[bi].Files
		var next []string

		for _, file := range files {
			if file != path {
				next = append(next, file)
			}
		}

		m.Plan.Branches[bi].Files = next
	}
}

func contains(values []string, value string) bool {
	for _, item := range values {
		if item == value {
			return true
		}
	}

	return false
}

func (m Model) View() string {
	if m.Width == 0 || m.Height == 0 {
		return "loading..."
	}

	if len(m.Repos) == 0 {
		return "no repositories configured"
	}

	topH := m.Height - 6
	if topH < 10 {
		topH = 10
	}

	header := m.renderHeader()
	var body string

	if m.Mode == enums.ModePush {
		body = m.renderPushMode(m.Width, topH)
	} else if m.Mode == enums.ModeQuickPush {
		body = m.renderQuickPushMode(m.Width, topH)
	} else if m.Mode == enums.ModeBranches {
		body = m.Branches.View(m.Width, topH)
	} else {
		body = m.renderGitMode(m.Width, topH)
	}

	footer := m.renderFooter(m.Width)

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m Model) renderQuickPushMode(w, h int) string {
	leftW := 54
	if w >= 170 {
		leftW = 66
	}

	rightW := w - leftW - 2
	if rightW < 60 {
		rightW = 60
	}

	left := m.renderQuickPushPanel(leftW, h)
	right := m.renderRepositoryGraph(rightW, h)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m Model) renderQuickPushPanel(w, h int) string {
	state := m.currentState()

	lines := []string{
		TitleStyle.Render("Quick Push"),
		MutedStyle.Render("current branch: " + state.Branch),
		MutedStyle.Render("operation: git add -A && git commit && git push"),
		"",
		TitleStyle.Render("Commit message"),
		InputStyle.Render(utils.TruncateLine(" "+m.QuickMessageInput+" ", w-4)),
		MutedStyle.Render("enter push · esc cancel"),
		"",
		TitleStyle.Render("Files included"),
	}

	if len(state.Files) == 0 {
		lines = append(lines, MutedStyle.Render("clean"))
	} else {
		for _, file := range state.Files {
			status := fmt.Sprintf("%s%s", normalizeStatus(file.Index), normalizeStatus(file.Worktree))
			line := fmt.Sprintf("%-2s %s", status, file.Path)
			lines = append(lines, OkStyle.Render(utils.TruncateLine("  "+line, w-4)))
		}
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)

	return FocusedBorder.Width(w).Height(h).Render(content)
}

func (m Model) renderGitMode(w, h int) string {
	if len(m.Repos) == 1 {
		return m.renderSingleRepoView(w, h)
	}

	leftW := responsiveLeftWidth(w)
	rightW := w - leftW - 2
	if rightW < 60 {
		rightW = 60
	}

	repoH := h / 2
	filesH := h - repoH

	left := lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderRepos(leftW, repoH),
		m.renderFiles(leftW, filesH),
	)

	right := m.renderDiff(rightW, h)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m Model) renderSingleRepoView(w, h int) string {
	filesW := responsiveLeftWidth(w)
	diffW := w - filesW - 2
	if diffW < 60 {
		diffW = 60
	}

	files := m.renderFiles(filesW, h)
	diff := m.renderDiff(diffW, h)

	return lipgloss.JoinHorizontal(lipgloss.Top, files, diff)
}

func (m Model) renderPushMode(w, h int) string {
	planW := 42
	fileW := 46

	if w >= 180 {
		planW = 52
		fileW = 58
	}

	graphW := w - planW - fileW - 3
	if graphW < 60 {
		graphW = 60
	}

	plan := m.renderPushPlans(planW, h)
	graph := m.renderRepositoryGraph(graphW, h)
	files := m.renderAssignmentFiles(fileW, h)

	return lipgloss.JoinHorizontal(lipgloss.Top, plan, graph, files)
}

func responsiveLeftWidth(w int) int {
	leftW := 38

	if w >= 150 {
		leftW = 44
	}

	if w >= 190 {
		leftW = 52
	}

	return leftW
}

func (m Model) renderHeader() string {
	repo := m.Repos[m.ActiveRepo]
	state := m.currentState()

	branch := state.Branch
	if branch == "" {
		branch = "no branch"
	}

	mode := "git"
	if m.Mode == enums.ModePush {
		mode = "push-plan"
	}
	if m.Mode == enums.ModeQuickPush {
		mode = "quick-push"
	}
	if m.Mode == enums.ModeBranches {
		mode = "branches"
	}

	width := m.Width
	if width < 20 {
		width = 20
	}

	textWidth := width - 4
	if textWidth < 10 {
		textWidth = 10
	}

	text := fmt.Sprintf(
		" GIT COCKPIT │ %s │ %s │ %s │ files:%d ",
		mode,
		repo.Name,
		branch,
		len(state.Files),
	)

	text = utils.TruncateLine(text, textWidth)

	return StatusBarStyle.
		Width(width).
		MaxWidth(width).
		Render(text)
}

func (m Model) renderRepos(w, h int) string {
	lines := []string{TitleStyle.Render("Repositories")}

	nameWidth := w - 12
	if nameWidth < 10 {
		nameWidth = 10
	}

	for i, repo := range m.Repos {
		state := m.States[repo.Path]

		marker := " "
		if len(state.Files) > 0 {
			marker = WarnStyle.Render("*")
		} else if state.Error == "" && state.Branch != "" {
			marker = OkStyle.Render("✓")
		}

		name := utils.TruncateLine(repo.Name, nameWidth)
		line := fmt.Sprintf("%s %-*s %3d", marker, nameWidth, name, len(state.Files))

		if i == m.ActiveRepo {
			line = SelectedStyle.Width(w - 4).Render(line)
		}

		lines = append(lines, line)
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)

	style := PanelBorder.Width(w).Height(h)
	if m.Focus == enums.PanelRepos {
		style = FocusedBorder.Width(w).Height(h)
	}

	return style.Render(content)
}

func (m Model) renderFiles(w, h int) string {
	state := m.currentState()
	lines := []string{TitleStyle.Render("Files")}

	if state.Error != "" {
		lines = append(lines, ErrorStyle.Render(state.Error))
	} else if len(state.Files) == 0 {
		lines = append(lines, MutedStyle.Render("clean"))
	} else {
		pathWidth := w - 12
		if pathWidth < 12 {
			pathWidth = 12
		}

		for i, file := range state.Files {
			status := fmt.Sprintf("%s%s", normalizeStatus(file.Index), normalizeStatus(file.Worktree))

			mark := " "
			if m.MarkedFiles[file.Path] {
				mark = "●"
			}

			path := utils.TruncateLine(file.Path, pathWidth)
			line := fmt.Sprintf("%s %-2s %-*s", mark, status, pathWidth, path)

			if i == m.ActiveFile {
				if m.MarkedFiles[file.Path] {
					line = MarkedStyle.Width(w - 4).Render(line)
				} else {
					line = SelectedStyle.Width(w - 4).Render(line)
				}
			}

			lines = append(lines, line)
		}
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)

	style := PanelBorder.Width(w).Height(h)
	if m.Focus == enums.PanelFiles {
		style = FocusedBorder.Width(w).Height(h)
	}

	return style.Render(content)
}

func (m Model) renderCurrentCommitLines(width int) []string {
	state := m.currentState()

	current := state.Branch
	if current == "" {
		current = "unknown"
	}

	lines := []string{
		MutedStyle.Render(utils.TruncateLine("base: "+current, width)),
	}

	if len(state.Commits) > 0 {
		commit := state.Commits[0]
		line := "head: " + commit.Hash

		if commit.Refs != "" {
			line += " " + commit.Refs
		}

		if commit.Message != "" {
			line += " " + commit.Message
		}

		lines = append(lines, MutedStyle.Render(utils.TruncateLine(line, width)))
	}

	return lines
}

func (m Model) renderPushPlans(w, h int) string {
	lines := []string{TitleStyle.Render("Push Plans")}
	lines = append(lines, m.renderCurrentCommitLines(w-4)...)
	lines = append(lines, "")

	if m.NamingBranch {
		label := "new"
		if m.QuickPush {
			label = "quick"
		}

		line := fmt.Sprintf("%s branch: %s", label, m.BranchInput)
		lines = append(lines, InputStyle.Render(utils.TruncateLine(" "+line+" ", w-4)))
		lines = append(lines, MutedStyle.Render("enter create · esc cancel"))
	}

	if m.NamingMessage {
		line := "message: " + m.MessageInput
		lines = append(lines, InputStyle.Render(utils.TruncateLine(" "+line+" ", w-4)))
		lines = append(lines, MutedStyle.Render("enter set · esc cancel"))
	}

	if len(m.Plan.Branches) == 0 {
		lines = append(lines, MutedStyle.Render("no push plans"))
		lines = append(lines, MutedStyle.Render("press b to create one"))
	} else {
		for i, branch := range m.Plan.Branches {
			style := planStyle(i)
			line := fmt.Sprintf("■ %s", branch.Name)

			if i == m.ActiveBranch {
				line = SelectedStyle.Width(w - 4).Render(line)
			} else {
				line = style.Render(utils.TruncateLine(line, w-4))
			}

			lines = append(lines, line)

			base := fmt.Sprintf("  base: %s", branch.Base)
			lines = append(lines, MutedStyle.Render(utils.TruncateLine(base, w-6)))

			status := branch.Status
			if status == "" {
				status = "pending"
			}

			lines = append(lines, MutedStyle.Render(utils.TruncateLine("  status: "+status, w-6)))

			message := branch.Message
			if strings.TrimSpace(message) == "" {
				message = "<missing commit message>"
				lines = append(lines, WarnStyle.Render(utils.TruncateLine("  msg: "+message, w-6)))
			} else {
				lines = append(lines, OkStyle.Render(utils.TruncateLine("  msg: "+message, w-6)))
			}

			if branch.Error != "" {
				lines = append(lines, ErrorStyle.Render(utils.TruncateLine("  error: "+branch.Error, w-6)))
			}

			if len(branch.Files) == 0 {
				lines = append(lines, MutedStyle.Render("  no files assigned"))
			} else {
				for _, file := range branch.Files {
					lines = append(lines, style.Render(utils.TruncateLine("  • "+file, w-6)))
				}
			}
		}
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)

	style := PanelBorder.Width(w).Height(h)
	if m.Focus == enums.PanelPlans {
		style = FocusedBorder.Width(w).Height(h)
	}

	return style.Render(content)
}

func (m Model) renderRepositoryGraph(w, h int) string {
	mode := "git graph"
	if m.GraphMode == enums.GraphText {
		mode = "branch text"
	}

	lines := []string{
		TitleStyle.Render("Repository Branches") + "  " + MutedStyle.Render("("+mode+", g toggles)"),
	}

	if m.GraphMode == enums.GraphText {
		lines = append(lines, m.renderBranchTextLines(w-4)...)
	} else {
		graph := m.currentState().Graph
		if strings.TrimSpace(graph) == "" {
			lines = append(lines, MutedStyle.Render("no graph available"))
		} else {
			for _, line := range strings.Split(graph, "\n") {
				lines = append(lines, renderGraphLine(utils.TruncateLine(line, w-4)))
			}
		}
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)

	style := PanelBorder.Width(w).Height(h)
	if m.Focus == enums.PanelGraph {
		style = FocusedBorder.Width(w).Height(h)
	}

	return style.Render(content)
}

func (m Model) renderBranchTextLines(width int) []string {
	state := m.currentState()

	current := state.Branch
	if current == "" {
		current = "current"
	}

	lines := []string{
		OkStyle.Render(utils.TruncateLine("current: "+current, width)),
		MutedStyle.Render(""),
		TitleStyle.Render("planned push branches"),
	}

	if len(m.Plan.Branches) == 0 {
		lines = append(lines, MutedStyle.Render("  none"))
		return lines
	}

	for i, branch := range m.Plan.Branches {
		style := planStyle(i)

		lines = append(lines, style.Render(utils.TruncateLine("  ├─ "+branch.Name, width)))
		lines = append(lines, MutedStyle.Render(utils.TruncateLine("  │  base: "+branch.Base, width)))

		filesText := fmt.Sprintf("  │  files: %d", len(branch.Files))
		lines = append(lines, MutedStyle.Render(utils.TruncateLine(filesText, width)))

		message := branch.Message
		if strings.TrimSpace(message) == "" {
			message = "<missing message>"
		}

		lines = append(lines, MutedStyle.Render(utils.TruncateLine("  │  msg: "+message, width)))
	}

	return lines
}

func renderGraphLine(line string) string {
	switch {
	case strings.Contains(line, "HEAD"):
		return OkStyle.Render(line)
	case strings.Contains(line, "origin/"):
		return WarnStyle.Render(line)
	case strings.Contains(line, "tag:"):
		return TitleStyle.Render(line)
	default:
		return MutedStyle.Render(line)
	}
}

func (m Model) renderAssignmentFiles(w, h int) string {
	state := m.currentState()

	lines := []string{
		TitleStyle.Render("Changed Files") + "  " + MutedStyle.Render("assigned files stay visible"),
	}

	if state.Error != "" {
		lines = append(lines, ErrorStyle.Render(state.Error))
	} else if len(state.Files) == 0 {
		lines = append(lines, MutedStyle.Render("clean"))
	} else {
		pathWidth := w - 10
		if pathWidth < 14 {
			pathWidth = 14
		}

		for i, file := range state.Files {
			status := fmt.Sprintf("%s%s", normalizeStatus(file.Index), normalizeStatus(file.Worktree))

			mark := " "
			if m.MarkedFiles[file.Path] {
				mark = "●"
			}

			branchIndex, branchName, assigned := m.assignedBranchForFile(file.Path)

			label := file.Path
			if assigned {
				label = fmt.Sprintf("%s → %s", file.Path, branchName)
			}

			line := fmt.Sprintf("%s %-2s %-*s", mark, status, pathWidth, utils.TruncateLine(label, pathWidth))

			if assigned {
				line = planStyle(branchIndex).Render(utils.TruncateLine(line, w-4))
			}

			if i == m.ActiveFile {
				if m.MarkedFiles[file.Path] {
					line = MarkedStyle.Width(w - 4).Render(line)
				} else {
					line = SelectedStyle.Width(w - 4).Render(line)
				}
			}

			lines = append(lines, line)
		}

		lines = append(lines, "")
		lines = append(lines, TitleStyle.Render("Legend"))

		for i, branch := range m.Plan.Branches {
			lines = append(lines, planStyle(i).Render(utils.TruncateLine("■ "+branch.Name, w-4)))
		}
	}

	content := utils.Fit(strings.Join(lines, "\n"), h-2)

	style := PanelBorder.Width(w).Height(h)
	if m.Focus == enums.PanelFiles {
		style = FocusedBorder.Width(w).Height(h)
	}

	return style.Render(content)
}

func (m Model) assignedBranchForFile(path string) (int, string, bool) {
	for i, branch := range m.Plan.Branches {
		for _, file := range branch.Files {
			if file == path {
				return i, branch.Name, true
			}
		}
	}

	return -1, "", false
}

func planStyle(index int) lipgloss.Style {
	switch index % 6 {
	case 0:
		return PlanColor1
	case 1:
		return PlanColor2
	case 2:
		return PlanColor3
	case 3:
		return PlanColor4
	case 4:
		return PlanColor5
	default:
		return PlanColor6
	}
}

func (m Model) renderDiff(w, h int) string {
	state := m.currentState()
	repo := m.Repos[m.ActiveRepo]
	mode := m.selectedDiffModeLabel()

	header := TitleStyle.Render("Diff") +
		"  " +
		MutedStyle.Render(repo.Name) +
		"  " +
		MutedStyle.Render(state.Branch) +
		"  " +
		MutedStyle.Render(mode)

	if state.Ahead > 0 || state.Behind > 0 {
		header += fmt.Sprintf("  ↑%d ↓%d", state.Ahead, state.Behind)
	}

	diff := state.Diff
	if diff == "" {
		diff = MutedStyle.Render("No diff selected.")
	} else if IsRealDiff(diff) {
		columnWidth := (w - 12) / 2
		if columnWidth >= 30 {
			diff = RenderSideBySideDiff(diff, columnWidth)
		} else {
			diff = RenderUnifiedDiff(diff, w-4)
		}
	} else {
		diff = MutedStyle.Render(utils.TruncateLine(diff, w-4))
	}

	content := utils.Fit(header+"\n\n"+diff, h-2)

	return PanelBorder.Width(w).Height(h).Render(content)
}

func (m Model) renderFooter(w int) string {
	if w < 20 {
		w = 20
	}

	var commands string

	if m.Mode == enums.ModeQuickPush {
		commands = "enter commit+push · esc cancel · type message"
	} else if m.Mode == enums.ModeBranches {
		commands = "enter checkout · tab commits/files · j/k move · r refresh · esc back · q"
	} else if m.Mode == enums.ModePush {
		commands = "esc git · g graph/text · C msg · e push · E all · tab panel · space mark · b plan · a assign · x unassign · r refresh · q"
	} else {
		commands = "p planner · P quick push · B branches · j/k move · h/l/tab panel · s/u stage · S/U all · r refresh · q"
	}

	status := m.StatusMsg

	if m.NamingBranch {
		if m.QuickPush {
			status = "quick push branch: " + m.BranchInput
		} else {
			status = "branch: " + m.BranchInput
		}
	}

	if m.NamingMessage {
		status = "message: " + m.MessageInput
	}

	if m.Mode == enums.ModeQuickPush {
		status = "quick message: " + m.QuickMessageInput
	}

	if status == "" {
		status = "ready"
	}

	lineWidth := w - 6
	if lineWidth < 10 {
		lineWidth = 10
	}

	commands = utils.TruncateLine(commands, lineWidth)
	status = utils.TruncateLine(status, lineWidth)

	return PanelBorder.
		Width(w - 2).
		MaxWidth(w - 2).
		Render(commands + "\n" + MutedStyle.Render(status))
}

func (m Model) currentState() models.GitState {
	if len(m.Repos) == 0 {
		return models.GitState{}
	}

	if m.ActiveRepo < 0 || m.ActiveRepo >= len(m.Repos) {
		return models.GitState{}
	}

	return m.States[m.Repos[m.ActiveRepo].Path]
}

func (m Model) selectedFile() (models.FileStatus, bool) {
	files := m.currentState().Files

	if len(files) == 0 || m.ActiveFile < 0 || m.ActiveFile >= len(files) {
		return models.FileStatus{}, false
	}

	return files[m.ActiveFile], true
}

func (m Model) selectedDiffModeLabel() string {
	file, ok := m.selectedFile()
	if !ok {
		return "diff"
	}

	if hasWorktreeChange(file) {
		return "unstaged diff"
	}

	if hasStagedChange(file) {
		return "staged diff"
	}

	return "diff"
}

func hasStagedChange(file models.FileStatus) bool {
	return file.Index != "" && file.Index != " " && file.Index != "." && file.Index != "?"
}

func hasWorktreeChange(file models.FileStatus) bool {
	return file.Worktree != "" && file.Worktree != " " && file.Worktree != "."
}

func (m Model) loadDiffCmd() tea.Cmd {
	if len(m.Repos) == 0 || m.ActiveRepo < 0 || m.ActiveRepo >= len(m.Repos) {
		return nil
	}

	return func() tea.Msg {
		repo := m.Repos[m.ActiveRepo]
		state := m.currentState()

		file, ok := m.selectedFile()
		if ok {
			state.Diff = m.Git.ReadDiffForFile(repo, file)
		} else {
			state.Diff = ""
		}

		return models.RepoLoadedMsg{Index: m.ActiveRepo, State: state}
	}
}

func loadRepoCmd(git interfaces.GitProvider, index int, repo models.Repo) tea.Cmd {
	return func() tea.Msg {
		state := git.ReadGitState(repo)

		if len(state.Files) > 0 {
			state.Diff = git.ReadDiffForFile(repo, state.Files[0])
		}

		return models.RepoLoadedMsg{Index: index, State: state}
	}
}

func (m Model) executeQuickPushCmd(message string) tea.Cmd {
	repo := m.Repos[m.ActiveRepo]

	return func() tea.Msg {
		err := m.Git.ExecuteQuickPush(repo, message)
		return models.QuickPushExecutedMsg{Err: err}
	}
}

func (m Model) executeSelectedPlanCmd() tea.Cmd {
	if m.Mode != enums.ModePush {
		return nil
	}

	if len(m.Plan.Branches) == 0 {
		return nil
	}

	if m.ActiveBranch < 0 || m.ActiveBranch >= len(m.Plan.Branches) {
		return nil
	}

	index := m.ActiveBranch
	plan := m.Plan.Branches[index]
	repo := m.Repos[m.ActiveRepo]

	return func() tea.Msg {
		err := m.Git.ExecutePushPlan(repo, models.PushPlanExecution{
			Branch:  plan.Name,
			Base:    plan.Base,
			Message: plan.Message,
			Files:   plan.Files,
		})

		return models.PlanExecutedMsg{Index: index, Err: err}
	}
}

func (m Model) executeAllPlansCmd() tea.Cmd {
	if m.Mode != enums.ModePush {
		return nil
	}

	if len(m.Plan.Branches) == 0 {
		return nil
	}

	plans := append([]interfaces.PlannedBranch(nil), m.Plan.Branches...)
	repo := m.Repos[m.ActiveRepo]

	return func() tea.Msg {
		var results []models.PlanExecutedMsg

		for index, plan := range plans {
			err := m.Git.ExecutePushPlan(repo, models.PushPlanExecution{
				Branch:  plan.Name,
				Base:    plan.Base,
				Message: plan.Message,
				Files:   plan.Files,
			})

			results = append(results, models.PlanExecutedMsg{
				Index: index,
				Err:   err,
			})

			if err != nil {
				break
			}
		}

		return models.AllPlansExecutedMsg{Results: results}
	}
}

func (m Model) stageSelectedCmd() tea.Cmd {
	file, ok := m.selectedFile()
	if !ok {
		return nil
	}

	repo := m.Repos[m.ActiveRepo]

	return runActionCmd(func() error {
		return m.Git.StageFile(repo, file.Path)
	}, "staged "+file.Path)
}

func (m Model) unstageSelectedCmd() tea.Cmd {
	file, ok := m.selectedFile()
	if !ok {
		return nil
	}

	repo := m.Repos[m.ActiveRepo]

	return runActionCmd(func() error {
		return m.Git.UnstageFile(repo, file.Path)
	}, "unstaged "+file.Path)
}

func (m Model) stageAllCmd() tea.Cmd {
	repo := m.Repos[m.ActiveRepo]

	return runActionCmd(func() error {
		return m.Git.StageAll(repo)
	}, "staged all")
}

func (m Model) unstageAllCmd() tea.Cmd {
	repo := m.Repos[m.ActiveRepo]

	return runActionCmd(func() error {
		return m.Git.UnstageAll(repo)
	}, "unstaged all")
}

func runActionCmd(fn func() error, message string) tea.Cmd {
	return func() tea.Msg {
		err := fn()
		return models.ActionDoneMsg{
			Message: message,
			Err:     err,
		}
	}
}

func normalizeStatus(s string) string {
	if s == " " || s == "" {
		return "."
	}

	return s
}
