package controllers

import "fmt"

func (m *Model) toggleMarkedFile() {
	file, ok := m.selectedFile()
	if !ok {
		return
	}

	key := file.Path
	if m.State.MarkedFiles[key] {
		delete(m.State.MarkedFiles, key)
		m.State.StatusMsg = "unmarked " + key
	} else {
		m.State.MarkedFiles[key] = true
		m.State.StatusMsg = "marked " + key
	}
}

func (m *Model) assignMarkedToBranch() {
	if len(m.State.Plan.Branches) == 0 {
		m.State.StatusMsg = "create a push plan first with b"
		return
	}

	if m.State.ActiveBranch < 0 || m.State.ActiveBranch >= len(m.State.Plan.Branches) {
		m.State.ActiveBranch = 0
	}

	branch := &m.State.Plan.Branches[m.State.ActiveBranch]
	count := 0

	for path := range m.State.MarkedFiles {
		m.removeFileFromAllBranches(path)

		if !contains(branch.Files, path) {
			branch.Files = append(branch.Files, path)
			count++
		}

		delete(m.State.MarkedFiles, path)
	}

	m.State.StatusMsg = fmt.Sprintf("assigned %d file(s) to %s", count, branch.Name)
}

func (m *Model) removeSelectedFileFromPlan() {
	file, ok := m.selectedFile()
	if !ok {
		return
	}

	m.removeFileFromAllBranches(file.Path)
	delete(m.State.MarkedFiles, file.Path)
	m.State.StatusMsg = "removed from push plan: " + file.Path
}

func (m *Model) removeFileFromAllBranches(path string) {
	for bi := range m.State.Plan.Branches {
		files := m.State.Plan.Branches[bi].Files
		var next []string

		for _, file := range files {
			if file != path {
				next = append(next, file)
			}
		}

		m.State.Plan.Branches[bi].Files = next
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
