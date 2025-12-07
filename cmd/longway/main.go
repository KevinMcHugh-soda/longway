package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

const (
	tickInterval = 250 * time.Millisecond
)

type model struct {
	progress progress.Model
	percent  float64
	stage    int
	stages   []string
	paused   bool
	width    int
	height   int
}

func newModel() model {
	p := progress.New(
		progress.WithDefaultGradient(),
	)

	return model{
		progress: p,
		stages: []string{
			"Plot your route through derelict sectors (FTL vibes).",
			"Draft your loadout and perks (Slay the Spire energy).",
			"Hit the notes to break through (Clone Hero / YARG).",
			"Watch the world react to your performance.",
		},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tick(), tea.EnterAltScreen)
}

func tick() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case " ":
			m.paused = !m.paused
		case "n":
			m.advanceStage()
		case "r":
			m.percent = 0
			m.paused = false
		}
	case tickMsg:
		if !m.paused {
			m.percent += 0.02
			if m.percent >= 1 {
				m.percent = 1
				m.paused = true
			}
		}
		return m, tick()
	}

	return m, nil
}

func (m *model) advanceStage() {
	if m.stage+1 < len(m.stages) {
		m.stage++
		m.percent = 0
		m.paused = false
	}
}

func (m model) View() string {
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F0D94A")).
		Bold(true).
		Underline(true).
		Render("Long Way To The Top")

	sub := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9AE6FF")).
		Render("A roguelike rhythm climb inspired by FTL, Slay the Spire, and Balatro.")

	stageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#B6EEA6")).
		Bold(true)

	body := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#6C7086")).
		Padding(1, 2).
		Width(max(60, m.width-4))

	progressBar := m.progress.ViewAs(m.percent)

	stageLabel := fmt.Sprintf("Stage %d/%d", m.stage+1, len(m.stages))
	stageDesc := m.stages[m.stage]

	status := "Press space to pause/resume. n: next stage, r: reset, q: quit."
	if m.paused && m.percent >= 1 {
		status = "Stage complete. Press n to chart the next stretch."
	} else if m.paused {
		status = "Paused. Press space to continue climbing."
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		stageStyle.Render(stageLabel),
		stageDesc,
		"",
		progressBar,
		"",
		status,
	)

	doc := lipgloss.JoinVertical(lipgloss.Left, title, sub, "", body.Render(content))

	if m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, doc)
	}

	return doc
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	m := newModel()
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
