package timer

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/gen2brain/beeep"
)

type TickMsg time.Time

type SessionType int

const (
	WorkSession SessionType = iota
	BreakSession
	LongBreakSession
)

const (
	WorkDuration      = 25 * time.Minute
	BreakDuration     = 5 * time.Minute
	LongBreakDuration = 25 * time.Minute

	BarWidth = 25

	FilledBlock = "█"
	EmptyBlock  = "░"
	StartSymbol = "▶"
	PauseSymbol = "⏸"
)

// Styling variables for lipgloss
var (
	backgroundColor = lipgloss.Color("#211e20")
	dimTextColor    = lipgloss.Color("#555568")
	warmTextColor   = lipgloss.Color("#a0a08b")
	coolTextColor   = lipgloss.Color("#e9efec")

	defaultStyle = lipgloss.NewStyle().Background(backgroundColor).Foreground(warmTextColor)
	titleStyle   = lipgloss.NewStyle().Foreground(dimTextColor).Bold(true).Inherit(defaultStyle)
	timeStyle    = lipgloss.NewStyle().Foreground(coolTextColor).Inherit(defaultStyle)
	helpStyle    = lipgloss.NewStyle().Foreground(dimTextColor).Faint(true).Inherit(defaultStyle)
)

type Model struct {
	Remaining   time.Duration
	Running     bool
	SessionType SessionType
	Pomodoros   int
	Width       int
	Height      int
}

// Helper Functions -----------------------------------------------------------

func New() Model {
	return Model{
		Remaining:   WorkDuration,
		SessionType: WorkSession,
	}
}

func sendTickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return TickMsg(t) })
}

func sessionDuration(st SessionType) time.Duration {
	switch st {
	case WorkSession:
		return WorkDuration
	case BreakSession:
		return BreakDuration
	case LongBreakSession:
		return LongBreakDuration
	default:
		return WorkDuration
	}
}

func nextSession(m Model) Model {
	if m.SessionType == WorkSession {
		// Every four sessions give a long break (after first session)
		if (m.Pomodoros%4 == 0) && (m.Pomodoros > 0) {
			m.SessionType = LongBreakSession
		} else {
			m.SessionType = BreakSession
		}
	} else {
		m.SessionType = WorkSession
	}
	m.Remaining = sessionDuration(m.SessionType)
	return m
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	m := int(d.Minutes())
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", m, s)
}

func notify(title, body string) tea.Cmd {
	return func() tea.Msg {
		_ = beeep.Alert(title, body, "")
		return nil
	}
}

func sessionString(st SessionType) string {
	switch st {
	case WorkSession:
		return "Work"
	case BreakSession:
		return "Break"
	case LongBreakSession:
		return "Long Break"
	default:
		return "Work"
	}
}

// Bubbletea Functions --------------------------------------------------------

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "space":
			m.Running = !m.Running
			if m.Running {
				return m, sendTickCmd()
			}
			return m, nil
		case "r":
			m.Running = false
			m.Remaining = sessionDuration(m.SessionType)
		case "s":
			m.Running = false
			m = nextSession(m)
		}

	case TickMsg:
		// Stops timer from updating if pause called in middle of tick
		if !m.Running {
			return m, nil
		}

		if m.Remaining > 0 {
			m.Remaining -= time.Second
			return m, sendTickCmd()
		}

		var alertCmd tea.Cmd
		if m.SessionType == WorkSession {
			m.Pomodoros++
		}
		body := fmt.Sprintf("%s Complete", sessionString(m.SessionType))
		alertCmd = notify("Pomodoro", body)

		m = nextSession(m)
		// Batch the notify to not block update loop
		return m, tea.Batch(sendTickCmd(), alertCmd)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}
	return m, nil
}

func (m Model) View() tea.View {
	curSession := sessionString(m.SessionType)

	// Handle the timer
	prettyTime := formatDuration(m.Remaining)

	// Handle the progress bar
	totalTime := sessionDuration(m.SessionType)
	elapsedTime := totalTime - m.Remaining
	progressFraction := float64(elapsedTime) / float64(totalTime)
	filledChars := int(progressFraction * float64(BarWidth))
	emptyChars := BarWidth - filledChars
	barString := strings.Repeat(FilledBlock, filledChars) + strings.Repeat(EmptyBlock, emptyChars)

	// Render with lipgloss styles
	title := titleStyle.Render("Pomodoro Timer")
	prettyTime = timeStyle.Render(prettyTime)
	barString = titleStyle.Render(barString)

	// If paused pause symbol is bright, play is dim: vice versa if running
	var runningString, pauseString string
	if m.Running {
		runningString = defaultStyle.Render(StartSymbol)
		pauseString = helpStyle.Render(PauseSymbol)
	} else {
		runningString = helpStyle.Render(StartSymbol)
		pauseString = defaultStyle.Render(PauseSymbol)
	}

	pomodoros := fmt.Sprintf("Pomodoros: %d", m.Pomodoros)

	curSession = defaultStyle.Render("Session: " + curSession)
	status := defaultStyle.Render(runningString + " " + pauseString)
	help := helpStyle.Render("[space] start/pause [s] skip [r] reset [q] quit")
	uiString := fmt.Sprintf("%s\n\n%s\n%s %s\n%s\n\n%s\n\n%s",
		title, curSession, barString, prettyTime, status, pomodoros, help)

	renderedContent := defaultStyle.Render(uiString)
	renderedContent = lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, renderedContent)

	view := tea.View{
		AltScreen: true,
		Content:   renderedContent,
	}

	return view
}
