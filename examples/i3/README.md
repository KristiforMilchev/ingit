# i3 Integration

This setup allows launching `ingit` directly from the currently focused terminal location using i3.

The workflow:

- Open a terminal
- `cd` into a repository
- Press a keybind
- `ingit` launches in that repository

This works especially well with:

- i3
- kitty
- alacritty
- wezterm
- foot
- xterm

---

# Overview

The integration consists of:

```text
~/.local/bin/ingit-here
~/.config/i3/config
```

The launcher script:

1. Detects the currently focused terminal window
2. Reads its current working directory
3. Opens `ingit` in that directory

---

# Step 1 — Build ingit

From the project:

```bash
go build -o ingit .
chmod +x ./ingit
```

Optional global install:

```bash
mkdir -p ~/.local/bin
cp ./ingit ~/.local/bin/ingit
chmod +x ~/.local/bin/ingit
```

Verify:

```bash
ingit
```

---

# Step 2 — Install Dependencies

Required:

```bash
sudo apt install xdotool
```

Optional but recommended:

```bash
sudo apt install jq
```

---

# Step 3 — Create Launcher Script

Create:

```bash
mkdir -p ~/.local/bin
nano ~/.local/bin/ingit-here
```

Paste:

```bash
#!/usr/bin/env bash

set -euo pipefail

WINDOW_ID=$(xdotool getwindowfocus)

PID=$(xprop -id "$WINDOW_ID" _NET_WM_PID \
  | awk '{print $3}')

if [ -z "${PID:-}" ]; then
    notify-send "ingit" "Unable to detect focused process"
    exit 1
fi

TARGET_PID="$PID"

while true; do
    CHILD=$(pgrep -P "$TARGET_PID" | head -n1 || true)

    if [ -z "$CHILD" ]; then
        break
    fi

    TARGET_PID="$CHILD"
done

CWD=$(readlink -f "/proc/$TARGET_PID/cwd" 2>/dev/null || true)

if [ -z "${CWD:-}" ]; then
    notify-send "ingit" "Unable to detect terminal directory"
    exit 1
fi

if [ ! -d "$CWD/.git" ]; then
    notify-send "ingit" "Not inside a git repository"
    exit 1
fi

exec kitty --directory "$CWD" ingit
```

Make executable:

```bash
chmod +x ~/.local/bin/ingit-here
```

---

# Step 4 — Add i3 Keybinding

Open:

```bash
nano ~/.config/i3/config
```

Add:

```text
bindsym $mod+Shift+g exec --no-startup-id ~/.local/bin/ingit-here
```

Reload i3:

```bash
i3-msg reload
```

or:

```text
Mod+Shift+R
```

---

# Usage

Example:

```bash
cd ~/Projects/project
```

Press:

```text
Mod+Shift+G
```

`ingit` opens directly inside that repository.

---

You can combine this with the VS Code task integration from `README_VSCODE.md`.

---

# Multi-Terminal Support

## Kitty

```bash
exec kitty --directory "$CWD" ingit
```

## Alacritty

```bash
exec alacritty --working-directory "$CWD" -e ingit
```

## WezTerm

```bash
exec wezterm start --cwd "$CWD" ingit
```

## Foot

```bash
exec foot --working-directory="$CWD" ingit
```

---
