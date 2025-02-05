package cookies

import (
	"fmt"
	"os"
	"regexp"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

const (
	authCookieIndex = iota
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

type model struct {
	inputs  []textinput.Model
	focused int
	err     error
}

func AskForCookie() string {
	// Ask for the cookie, showing how to obtain it
	p := tea.NewProgram(cookieModel())
	model_out, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting app: %v\n", err)
		os.Exit(1)
	}

	return model_out.(model).inputs[authCookieIndex].Value()

}

func cookieValidator(s string) error {
	// Token should be a string of 22 characters, that matches the regular expression
	if s != "" && regexp.MustCompile(`[a-zA-Z0-9]{25,}`).MatchString(s) && len(s) > 25 {
		return nil
	}
	return fmt.Errorf("cookie is invalid")
}

func cookieModel() model {
	var inputs []textinput.Model = make([]textinput.Model, 1)
	focusSet := false
	inputs[authCookieIndex] = textinput.New()
	inputs[authCookieIndex].Placeholder = "AulaGlobal auth cookie"
	inputs[authCookieIndex].CharLimit = 32
	inputs[authCookieIndex].Width = 32
	inputs[authCookieIndex].Prompt = ""
	inputs[authCookieIndex].Validate = cookieValidator
	if !focusSet {
		inputs[authCookieIndex].Focus()
		focusSet = true
	}

	return model{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				return m, tea.Quit
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc, tea.KeyCtrlD:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		default:

		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	// TODO: use colors 'n stuff
	return fmt.Sprintf(`
Token file not found.

You must provide a cookie to obtain the token. To do this:
1. Log into Aula Global
2. Open the developer tools (F12)
3. Go to the console tab, and run the following command:

   console.log(('; ' + document.cookie).split('; MoodleSessionag=').pop().split(';').shift())

 %s
 %s

 %s
`,
		inputStyle.Width(30).Render("Cookie"),
		m.inputs[authCookieIndex].View(),
		continueStyle.Render("Continue (enter)->"),
	) + "\n"
}

// nextInput focuses the next input field
func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *model) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
