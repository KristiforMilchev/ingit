# ingit

A terminal-native Git workspace manager designed for modern development workflows.

`ingit` focuses on fast repository navigation, structured commit planning, multi-repository management, and IDE integration — all without leaving the terminal.

Yes, tools like :contentReference[oaicite:0]{index=0} already exist, but `ingit` approaches the problem differently:
- editor agnostic
- terminal first
- mono-repository aware
- optimized for fast iterative commits
- designed to integrate into existing workflows instead of replacing them

---

## Preview

> Add screenshots here once available.

```md
![Overview](./assets/overview.png)
![Push Planner](./assets/planner.png)
```

---

# Why ingit

Git CLI remains extremely powerful, especially for advanced operations and recovery workflows.

However, modern software development has evolved significantly:
- large mono repositories
- multiple active services
- rapid iterative commits
- feature branch explosion
- staged execution plans
- terminal-driven workflows
- tiled window manager setups

In these environments, repeatedly typing Git commands, switching terminals, and manually organizing changes can become unnecessarily disruptive.

`ingit` aims to streamline that process.

It provides a structured TUI workflow for:
- managing multiple repositories simultaneously
- organizing changes into execution plans
- rapidly creating commits
- exploring branches visually
- operating entirely from the keyboard
- integrating directly into editors and window managers

---

# Features

## Multi-Repository Overview

Designed for mono-repository and multi-service environments.

Track and manage changes across multiple repositories from a single interface without switching terminals or editor windows.

---

## Quick Push Workflow

Need to commit and push immediately to the current branch?

`ingit` provides a dedicated quick-push flow optimized for minimal interaction.

---

## Push Planner

One of the core features of `ingit`.

Large unorganized changes can be:
- split into multiple execution plans
- assigned to separate branches
- grouped by file ownership
- committed independently
- executed sequentially

All without mutating repository state until execution.

This allows safe commit restructuring even after development work has already happened.

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

`ingit` automatically simplifies the UI and removes unnecessary repository management panels.

---

## Terminal Native

Built specifically for:
- keyboard-driven workflows
- tiling window managers
- IDE terminal integration
- remote development
- SSH workflows

---

# Integrations

`ingit` is intentionally editor agnostic.

Anything capable of launching a terminal can integrate with it.

Ready-made integrations are included for:
- VS Code
- VSCodium
- i3

Additional integrations are welcome through pull requests.

---

# Building From Source

## Requirements

- Go 1.24+ recommended

## Build

```bash
git clone <repo>
cd ingit

go mod tidy
go build -o ingit .
```

Optional global install:

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

# Integration Examples

See:
- `README_VSCODE.md`
- `README_I3.md`

for ready-made integration setups.

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
