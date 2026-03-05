# Pomodoro Timer

<img width="426" height="179" alt="image" src="https://github.com/user-attachments/assets/a21f0d19-c158-47b4-827b-4084213b7476" />


A terminal-based Pomodoro timer built with Go, using
[Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss)


## Features

- **Pomodoro cycling** -- 25-minute work sessions, 5-minute breaks, and
  25-minute long breaks every four pomodoros
- **Progress bar**
- **System notifications** when each session completes

## Getting Started

```
go build -o pomodoro ./cmd/pomodoro
./pomodoro
```

Or run directly:

```
go run ./cmd/pomodoro
```

## Controls

| Key       | Action                       |
|-----------|------------------------------|
| `space`   | Start / pause the timer      |
| `s`       | Skip to the next session     |
| `r`       | Reset the current session    |
| `q`       | Quit                         |

## Requirements

- Go 1.24+
- A terminal with unicode support
