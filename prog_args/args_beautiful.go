package prog_args

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

const (
	dirIota = iota
	corIota
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

// Validator functions to ensure valid input
func corValidator(s string) error {
	// Number of cores should be an integer
	_, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("the number of cores must be an integer")
	}
	return nil
}

func dirValidator(s string) error {
	// The directory should be a string of less than 40 characters
	if len(s) > 40 {
		return fmt.Errorf("directory is too long")
	}

	return nil
}

func initialModel(dirStr *string, cores int) model {
	var inputs []textinput.Model = make([]textinput.Model, 2)
	focusAlreadySet := false
	focusedResult := 0

	// Directory input setup
	inputs[dirIota] = textinput.New()
	inputs[dirIota].Placeholder = "(current directory)"
	inputs[dirIota].CharLimit = 40
	inputs[dirIota].Width = 30
	inputs[dirIota].Prompt = ""
	if *dirStr != "" {
		inputs[dirIota].SetValue(*dirStr)
	} else {
		inputs[dirIota].Focus()
		focusAlreadySet = true
		focusedResult = 0
	}
	inputs[dirIota].Validate = dirValidator

	// Cores input setup
	inputs[corIota] = textinput.New()
	if cores == 0 {
		cores = runtime.NumCPU() / 2 // half of the total CPUs
	}
	if !focusAlreadySet {
		inputs[corIota].Focus()
		focusedResult = 1
	}
	inputs[corIota].SetValue(strconv.Itoa(cores))
	inputs[corIota].Placeholder = "Number of cores to use"
	inputs[corIota].CharLimit = 3
	inputs[corIota].Width = 5
	inputs[corIota].Prompt = ""
	inputs[corIota].Validate = corValidator

	return model{
		inputs:  inputs,
		focused: focusedResult,
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
	return fmt.Sprintf(`Input the save directory and number of cores to use (-1 means all cores):

 %s   %s
 %s   %s

 %s
`,
		inputStyle.Width(30).Render("Directory"),
		inputStyle.Width(5).Render("Cores"),
		m.inputs[dirIota].View(),
		m.inputs[corIota].View(),
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
