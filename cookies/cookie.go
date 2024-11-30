package prog_args

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/Astrak00/AGDownloader/types"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/fatih/color"
)

type (
	errMsg error
)

const (
	authCookieIndex = iota
	dirIota
	corIota
)

const (
	hotPink          = lipgloss.Color("#FF06B7")
	darkGray         = lipgloss.Color("#767676")
	CookieText       = "To do this, the easiest way is to log in to Aula Global, open the developer tools, go to the console tab and run the following command:"
	ObtainCookieText = `function getCookie() {
		const value = '; ' + document.cookie;
		const parts = value.split('; MoodleSessionag=');
		if (parts.length === 2) console.log(parts.pop().split(';').shift());
	  } 
	  getCookie();`
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

func GetTokenFromCookie(arguments types.Prog_args) string {
	p := tea.NewProgram(initialModel(&arguments.DirPath, arguments.MaxGoroutines))
	model_out, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting app: %v\n", err)
		os.Exit(1)
	}

	cookie := model_out.(model).inputs[authCookieIndex].Value()
	// Convert from cookie to token

	token := cookie_to_token(cookie)

	return token

}

// Validator functions to ensure valid input
func corValidator(s string) error {
	// Number of cores should be an integer
	_, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("cores must be an integer")
	}
	return nil
}

func tokenValidator(s string) error {
	// Token should be a string of 22 characters, that matches the regular expression
	if s != "" && regexp.MustCompile(`[a-zA-Z0-9]{25,}`).MatchString(s) && len(s) > 25 {
		return nil
	}
	return fmt.Errorf("cookie is invalid")
}

func dirValidator(s string) error {
	// The directory should be a string of less than 40 characters
	if len(s) > 40 {
		return fmt.Errorf("directory is too long")
	}

	return nil
}

func initialModel(dirStr *string, cores int) model {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	focusSet := false
	inputs[authCookieIndex] = textinput.New()
	inputs[authCookieIndex].Placeholder = "Authentication cookie from AulaGlobal website"
	inputs[authCookieIndex].CharLimit = 32
	inputs[authCookieIndex].Width = 32
	inputs[authCookieIndex].Prompt = ""
	inputs[authCookieIndex].Validate = tokenValidator
	inputs[authCookieIndex].Focus()

	inputs[dirIota] = textinput.New()
	inputs[dirIota].Placeholder = "downloaded_files"
	inputs[dirIota].CharLimit = 40
	inputs[dirIota].Width = 30
	inputs[dirIota].Prompt = ""
	if *dirStr != "" {
		inputs[dirIota].SetValue(*dirStr)
	} else if !focusSet {
		inputs[dirIota].Focus()
		focusSet = true
	}
	inputs[dirIota].Validate = dirValidator

	inputs[corIota] = textinput.New()
	if cores == -1 {
		cores = 1
	}
	if !focusSet {
		inputs[corIota].Focus()
	}
	inputs[corIota].SetValue(strconv.Itoa(cores))
	inputs[corIota].Placeholder = "Number of cores to use"
	inputs[corIota].CharLimit = 3
	inputs[corIota].Width = 5
	inputs[corIota].Prompt = ""
	inputs[corIota].Validate = corValidator

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
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		case tea.KeyCtrlD:
			// Exit the program if the user presses Ctrl+D
			color.Red("Press CTRL+C to exit")
			return m, tea.Quit
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
	return fmt.Sprintf(`Input the directory, token and number of cores to use:

 %s
 %s

 %s   %s
 %s   %s

 %s
`,
		inputStyle.Width(30).Render("Cookie"),
		m.inputs[authCookieIndex].View(),
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