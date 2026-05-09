# ingit

A terminal-native Git workspace manager designed for modern development workflows.

`ingit` focuses on fast repository navigation, structured commit planning, multi-repository management, branch exploration, and IDE integration — all without leaving the terminal.

Built for:
- keyboard-first developers
- mono repository environments
- terminal heavy workflows
- tiled window manager users
- fast iterative commits

---

# Preview

## Multi Repository Workflow

![Multi Repository Overview](./assets/multi-repo-overview.png)

---

## Side-by-Side Diff Viewer

![Side By Side Diff](./assets/side-by-side-diff-view.png)

---

## Push Planner

### Create Branch Plans

![Planner Branch Creation](./assets/planner-branch-creation.png)

### Assign Files To Plans

![Planner Assignment](./assets/planner-file-assignment.png)

### Graph View

![Planner Graph](./assets/planner-graph-view.png)

### Execute Plans

![Planner Execution](./assets/planner-execution-complete.png)

---

## Branch Explorer

### Branch Overview

![Branches](./assets/branches-view.png)

### Explore Branch Contents

![Branch File Explorer](./assets/branch-file-explorer.png)

---

# Why ingit

Git CLI remains extremely powerful, especially for:
- recovery workflows
- advanced history manipulation
- low-level repository operations

However modern development workflows evolved significantly beyond simple repository management.

Today developers commonly work with:
- multiple repositories
- large mono repositories
- many feature branches
- continuous iterative commits
- terminal based workflows
- editor integrated terminals

In these environments repeatedly switching terminals, editors, and manually organizing changes becomes unnecessary friction.

`ingit` attempts to reduce that friction while keeping the flexibility and transparency of Git itself.

---

# Features

## Multi Repository Overview

Designed specifically for mono repositories and multi-service environments.

Track changes across multiple repositories from a single interface.

No terminal switching.
No editor switching.
No repeated git status calls.

---

## Push Planner

One of the core features of `ingit`.

Large unorganized changes can be:
- split into multiple execution plans
- assigned to separate branches
- grouped by file ownership
- committed independently
- executed sequentially

without mutating repository state until execution.

This allows safe commit restructuring even after development work has already happened.

---

## Quick Push Workflow

Need to quickly commit and push into the current branch?

`ingit` provides a dedicated quick-push flow optimized for minimal interaction.

---

## Side-by-Side Diff Viewer

Built-in terminal diff visualization with:
- staged changes
- unstaged changes
- side-by-side rendering
- unified rendering fallback

---

## Branch Explorer

Browse:
- local branches
- remote branches
- commit history
- repository contents

and switch branches directly from the TUI.

---

## Context-Aware Layout

Single repository?

`ingit` automatically simplifies the layout and removes unnecessary repository management panels.

---

## Terminal Native

Built specifically for:
- keyboard-driven workflows
- SSH workflows
- remote development
- tiled window managers
- integrated IDE terminals

---

# Keyboard Driven Workflow

`ingit` is intentionally designed around fast keyboard navigation.

Core controls include:

| Key | Action |
|---|---|
| `j/k` | Navigate |
| `h/l` | Switch panels |
| `tab` | Cycle panels |
| `space` | Mark file |
| `b` | Create push plan |
| `a` | Assign files |
| `x` | Remove assigned file |
| `C` | Set commit message |
| `e` | Execute selected plan |
| `E` | Execute all plans |
| `P` | Quick push |
| `s/u` | Stage / unstage |
| `r` | Refresh |
| `g` | Toggle graph mode |

---

# Integrations

`ingit` is intentionally editor agnostic.

Anything capable of launching a terminal can integrate with it.

Ready-made examples are included for:
- VS Code
- VSCodium
- i3

See:
- `examples/code/`
- `examples/i3/`

---

# Building From Source

## Requirements

- Go 1.24+ recommended

---

## Clone

```bash
git clone <repo>
cd ingit
```

---

## Build

```bash
go mod tidy
go build -o ingit .
```

---

## Optional Global Install

```bash
mkdir -p ~/.local/bin
cp ./ingit ~/.local/bin/ingit
chmod +x ~/.local/bin/ingit
```

Verify installation:

```bash
ingit
```

---

# Example Usage

Open inside current repository:

```bash
ingit
```

Open specific repository:

```bash
ingit ~/Projects/my-project
```

---

# VS Code / VSCodium Integration

Example integration files are included under:

```text
examples/code/
```

You can configure:
- custom keybindings
- dedicated tasks
- fullscreen terminal launches
- workspace aware launches

---

# i3 Integration

`ingit` was designed to integrate naturally into tiled window manager workflows.

Example i3 configurations are included under:

```text
examples/i3/
```

This allows:
- launching from active terminal paths
- workspace aware execution
- floating or tiled layouts
- fast keyboard-driven repository management

---

# Philosophy

`ingit` does not attempt to replace Git.

It attempts to reduce workflow friction around Git.

The goal is to make:
- frequent commits
- branch organization
- multi-repository work
- execution planning

feel lightweight and immediate.

---

# Current State

`ingit` is actively evolving.

Features and workflows are still expanding rapidly.

The project currently focuses on:
- workflow speed
- keyboard ergonomics
- repository visualization
- safe execution planning
- editor integration

---

# Contributing

Feature suggestions and pull requests are welcome.

If an integration for your editor or environment does not exist yet, feel free to open a PR.

The project aims to remain:
- lightweight
- terminal-native
- dependency minimal
- workflow focused
- editor agnostic

---

# License

MIT
