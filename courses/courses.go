package courses

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"

	types "github.com/Astrak00/AGDownloader/types"
)

// GetCourses obtains the courses, the localized name and ID, given a userID
// Returns a slice of courses
func GetCourses(token string, userID string, language int) (types.Courses, error) {
	fmt.Println("Fetching courses from AulaGlobal...")

	url := fmt.Sprintf(
		"https://%s%s?wstoken=%s&wsfunction=core_enrol_get_users_courses&userid=%s&moodlewsrestformat=json",
		types.Domain,
		types.Webservice,
		token,
		userID,
	)

	jsonData := types.GetJson(url)

	// Parse the json
	var userParsed types.WebUser
	err := json.Unmarshal(jsonData, &userParsed)
	if err != nil {
		log.Fatal(err)
	}

	// Get the names and IDs of the courses
	courses := make([]types.Course, 0, len(userParsed))
	for _, course := range userParsed {
		courseName := extractCourseNameByLanguage(course.Fullname, language)
		if !containsInvalidNames(courseName) {
			courses = append(courses, types.Course{Name: courseName, ID: strconv.Itoa(course.ID)})
		}
	}

	defer color.Green("Number of courses found: %d\n", len(courses))
	return courses, nil
}

// Check if the name of the course contains invalid names that should not be downloaded
func containsInvalidNames(name string) bool {
	invalidCourseNames := []string{
		"Convenio", "Delegación", "Secretaría",
		"Student Room", "Sala de Estudiantes", "Bachelor",
	}

	for _, invalidName := range invalidCourseNames {
		if strings.Contains(name, invalidName) {
			return true
		}
	}
	return false
}

// Get the names of the courses in Spanish and English
// This function localizes the names of the courses in Spanish and English
func extractCourseNameByLanguage(name string, lang int) string {
	// Define the first group of separators with priority.
	firstGroup := []string{"-1C", "-2C", "-1S", "-2S"}

	// Iterate over the first group to find the earliest separator.
	for _, sep := range firstGroup {
		if idx := strings.Index(name, sep); idx > 0 {
			if lang == 1 {
				return name[:idx+len(sep)]
			}
			return name[idx+len(sep):]
		}
	}

	// Define the second group of separators.
	secondGroup := []string{"Bachelor", "Student", "Convenio-Bilateral s"}

	// Iterate over the second group to find the earliest separator.
	for _, sep := range secondGroup {
		if idx := strings.Index(name, sep); idx != -1 {
			if lang == 1 {
				return name[:idx]
			}
			return name[idx:]
		}
	}

	// If no separators are found, return the original name.
	return name
}

// SelectCoursesInteractive is the entry point for prompting the user:
func SelectCoursesInteractive(language int, selectedCourses []string, courses types.Courses) []types.Course {
	if len(selectedCourses) != 0 && selectedCourses[0] == "all" {
		return courses
	} else if len(selectedCourses) == 0 {
		// Use our Bubble Tea-based checkboxes
		prompt := "Select the courses you want to download\n"
		coursesName := courses.GetCoursesName()

		selectedCourses = checkboxesCourses(prompt, coursesName)
	}

	coursesToDownload := make([]types.Course, 0, len(selectedCourses))
	courseMap := make(map[string]types.Course)
	for _, c := range courses {
		courseMap[c.Name] = c
	}

	for _, courseName := range selectedCourses {
		if course, exists := courseMap[courseName]; exists {
			coursesToDownload = append(coursesToDownload, course)
		}
	}
	return coursesToDownload
}

// checkboxesCourses uses Bubble Tea to allow the user to interactively
// select items by pressing up/down to move and space to toggle selection.
func checkboxesCourses(label string, opts []string) []string {
	m := initialModel(label, opts)

	// Run the Bubble Tea program
	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		fmt.Println("Error running Bubble Tea program:", err)
		os.Exit(1)
	}

	// Extract the selected items
	if mFinal, ok := finalModel.(model); ok {
		return mFinal.selectedItems()
	}
	return nil
}

// ----------------------------------------------------
// Below is a minimal Bubble Tea model for multi-select
// ----------------------------------------------------
type model struct {
	label    string
	cursor   int          // which item is currently highlighted
	items    []string     // all course names
	selected map[int]bool // track selected items by index
	done     bool         // signals we've pressed Enter
	viewport viewport.Model
	keymap   keymap
}

// Define key bindings we care about
type keymap struct {
	Up    key.Binding
	Down  key.Binding
	Space key.Binding
	Enter key.Binding
	Quit  key.Binding
	All   key.Binding
	None  key.Binding
}

// initialModel sets up the model with defaults
func initialModel(label string, items []string) model {
	m := model{
		label:    label,
		items:    items,
		selected: make(map[int]bool),
		cursor:   0,
		keymap: keymap{
			Up: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "move up"),
			),
			Down: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "move down"),
			),
			Space: key.NewBinding(
				key.WithKeys(" "),
				key.WithHelp("space", "toggle selection"),
			),
			Enter: key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "confirm selection"),
			),
			Quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q/ctrl+c", "quit"),
			),
			All: key.NewBinding(
				key.WithKeys("*", "right"),
				key.WithHelp("*", "select all"),
			),
			None: key.NewBinding(
				key.WithKeys("0", "left"),
				key.WithHelp("0", "select none"),
			),
		},
	}

	m.viewport = viewport.New(0, 0)
	return m
}

// Init is called when the program starts. We don't need to do anything here.
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages (keypresses, window size changes, etc.)
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		// Move cursor up
		case key.Matches(msg, m.keymap.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		// Move cursor down
		case key.Matches(msg, m.keymap.Down):
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		// Toggle selection
		case key.Matches(msg, m.keymap.Space):
			m.selected[m.cursor] = !m.selected[m.cursor]
		// Confirm (Enter) -> exit
		case key.Matches(msg, m.keymap.Enter):
			m.done = true
			return m, tea.Quit
		// Quit
		case key.Matches(msg, m.keymap.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keymap.All):
			for i := range m.items {
				m.selected[i] = true
			}
		case key.Matches(msg, m.keymap.None):
			for i := range m.items {
				m.selected[i] = false
			}
		}

	case tea.WindowSizeMsg:
		// If the window resizes, update viewport size
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
	}

	return m, nil
}

// View renders the UI each time Update is called.
func (m model) View() string {
	if m.done {
		// Once done, just return. Program will quit, returning to checkboxesCourses.
		return ""
	}

	s := m.label + "\n"

	for i, choice := range m.items {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // highlight the current line
		}

		checked := " " // not selected
		if m.selected[i] {
			checked = "x"
		}

		// [ ] or [x], plus cursor arrow, plus the course name
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	s += "\n(↑/↓ or k/j to navigate, space to toggle, enter to confirm, q to quit)\n(*/→ to select all, ←/0 to select none)"
	return s
}

// selectedItems returns the items that the user marked as selected
func (m model) selectedItems() []string {
	results := []string{}
	for i, selected := range m.selected {
		if selected {
			results = append(results, m.items[i])
		}
	}
	return results
}
