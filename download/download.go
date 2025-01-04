package download

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	types "github.com/Astrak00/AGDownloader/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	totalFiles     int32
	courseIDMap    map[string]string
	completedFiles int32
	currentFile    string
	errs           []string
}

func (m model) Init() tea.Cmd {
	// No initialization
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case progressMsg:
		atomic.AddInt32(&m.completedFiles, 1)
		m.currentFile = msg.fileName
		return m, nil

	case errorMsg:
		m.errs = append(m.errs, fmt.Sprintf("Error downloading %s: %v", msg.fileName, msg.err))
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	// Gradient styling

	progress := float64(m.completedFiles) / float64(m.totalFiles) * 100
	barWidth := 30
	filled := int(progress / 100 * float64(barWidth))
	empty := barWidth - filled

	filledBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Render(string(repeat('█', filled)))

	emptyBar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#353C49")).
		Render(string(repeat(' ', empty)))

	bar := fmt.Sprintf("%s%s %.1f%%", filledBar, emptyBar, progress)

	view := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Render(fmt.Sprintf("Downloading files...\n%s\nCompleted: %d/%d\n", bar, m.completedFiles, m.totalFiles))

	if m.currentFile != "" {
		courseName := m.courseIDMap[m.currentFile[:6]]
		if idx := strings.Index(courseName, "/"); idx != -1 {
			courseName = courseName[:idx-3]
		}
		view += lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF8C00")).
			Render(fmt.Sprintf("\nCurrent file: %s/%s\n", courseName, m.currentFile[7:]))
	}
	if len(m.errs) > 0 {
		view += lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF3333")).
			Render("\nErrors:\n")
		for _, e := range m.errs {
			view += fmt.Sprintf("- %s\n", e)
		}
	}
	view += "\nPress 'q' to quit.\n"
	return view
}

type progressMsg struct {
	fileName string
}

type errorMsg struct {
	fileName string
	err      error
}

func repeat(char rune, count int) []rune {
	out := make([]rune, count)
	for i := range out {
		out[i] = char
	}
	return out
}

// DownloadFiles orchestrates the file downloads and displays progress using Bubble Tea.
func DownloadFiles(filesStoreChan <-chan types.FileStore, maxGoroutines int, courses []types.Course) {
	totalFiles := len(filesStoreChan)

	// Convert the courses to a map for easy lookup
	courseIDMap := make(map[string]string)
	for _, course := range courses {
		courseIDMap[course.ID] = course.Name
	}

	m := model{
		totalFiles:  int32(totalFiles),
		courseIDMap: courseIDMap,
	}

	// Create the Bubble Tea program
	p := tea.NewProgram(m)

	// Start the program in a goroutine
	go func() {
		var wg sync.WaitGroup
		semaphore := make(chan struct{}, maxGoroutines)

		for fileStore := range filesStoreChan {
			wg.Add(1)
			go func(fileStore types.FileStore) {
				defer wg.Done()
				semaphore <- struct{}{}
				defer func() { <-semaphore }()
				if err := downloadFile(fileStore); err != nil {
					p.Send(errorMsg{fileName: fileStore.FileName, err: err})
				} else {
					p.Send(progressMsg{fileName: fileStore.FileName})
				}
			}(fileStore)
		}
		wg.Wait()

		// Quit the program after all downloads are complete
		p.Send(tea.Quit())
	}()

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error: %v\n", err)
	}
}

func downloadFile(fileStore types.FileStore) error {
	resp, err := http.Get(fileStore.FileURL)
	if err != nil {
		return fmt.Errorf("error downloading the file: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing body")
		}
	}(resp.Body)

	dir := filepath.Dir(fileStore.Dir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating the directory: %v", err)
	}

	out, err := os.Create(fileStore.Dir)
	if err != nil {
		return fmt.Errorf("error creating the file: %v", err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			fmt.Println("Error closing file")
		}
	}(out)

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error copying the file: %v", err)
	}

	return nil
}
