package prog_args

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model_param struct {
	cursor    int
	choices   []string
	selected  bool
	parameter string
}

func initialModel_token(argument string) model_param {
	return model_param{
		cursor:    0,
		choices:   []string{"Yes", "No"},
		selected:  false,
		parameter: argument,
	}
}

func (m model_param) Init() tea.Cmd {
	return nil
}

func (m model_param) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = true
			return m, tea.Quit
		case "q", "esc", "ctrl+c", "ctrl+d":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model_param) View() string {
	if m.selected {
		return fmt.Sprintf("You selected: %s to knowing %s\n", m.choices[m.cursor], m.parameter)
	}

	s := fmt.Sprintf("Do you know your %s?\n\n", m.parameter)

	for i, choice := range m.choices {
		cursor := " " // No cursor
		if m.cursor == i {
			cursor = ">" // Cursor
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\n(Use ↑/↓ or k/j to navigate, Enter to select, q to quit)"
	return s
}

func Knowledge_element(parameter string) bool {
	p := tea.NewProgram(initialModel_token(parameter))
	model, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting app: %v\n", err)
		os.Exit(1)
	}
	return model.(model_param).cursor == 0
}
