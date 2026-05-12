# VS Code Integration

`ingit` can be launched from VS Code as a dedicated terminal window using a VS Code task and keybinding.

This setup opens `ingit` in a new VS Code window/context and runs it inside the project folder.

## Requirements

- `ingit` binary must be built
- The binary must be executable
- VS Code or VSCodium must be installed
- The project must be opened as a workspace/folder

## Build

From the `ingit` project:

```bash
go build -o ingit .
chmod +x ./ingit
```

```bash
mkdir -p ~/.local/bin
cp ./ingit ~/.local/bin/ingit
chmod +x ~/.local/bin/ingit
```

Make sure this works:

```bash
ingit
```

## VS Code Task

Create:

```text
.vscode/tasks.json
```

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "ingit",
      "type": "shell",
      "command": "ingit",
      "options": {
        "cwd": "${workspaceFolder}"
      },
      "presentation": {
        "echo": false,
        "reveal": "always",
        "focus": true,
        "panel": "new",
        "clear": true
      },
      "problemMatcher": []
    }
  ]
}
```

If you do not install `ingit` globally, use the absolute path instead:

```json
"command": "/home/User/Projects/ingit/ingit"
```

## Keybinding

Open VS Code keybindings JSON:

```bash
code ~/.config/Code/User/keybindings.json
```

For VSCodium:

```bash
codium ~/.config/VSCodium/User/keybindings.json
```

Add:

```json
[
  {
    "key": "ctrl+shift+g g",
    "command": "workbench.action.tasks.runTask",
    "args": "ingit"
  }
]
```

Now pressing:

```text
Ctrl+Shift+G, then G
```

will launch `ingit`.

## Launch in a New VS Code Window

Create this script:

```bash
mkdir -p ~/.local/bin
nano ~/.local/bin/ingit-code
```

```bash
#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="${1:-$PWD}"

code --new-window "$PROJECT_DIR"

sleep 0.8

code --reuse-window "$PROJECT_DIR" \
  --command "workbench.action.tasks.runTask"
```

Make it executable:

```bash
chmod +x ~/.local/bin/ingit-code
```

Then run:

```bash
ingit-code /path/to/project
```
