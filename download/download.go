package download

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	errorlog "github.com/Astrak00/AGDownloader/errorlog"
	types "github.com/Astrak00/AGDownloader/types"
	"github.com/fatih/color"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	maxRetries     = 3
	initialBackoff = 1 * time.Second
)

type model struct {
	totalFiles     int32
	courseIDMap    map[string]string
	completedFiles int32
	currentFile    string
	errs           []string
	cancelled      bool
	errorLogger    *errorlog.ErrorLogger
	failedFiles    []types.FileStore
	failedFilesMux sync.Mutex
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
		errStr := fmt.Sprintf("Error downloading %s: %v", msg.fileName, msg.err)
		m.errs = append(m.errs, errStr)

		// Add to failed files list for retry
		m.failedFilesMux.Lock()
		m.failedFiles = append(m.failedFiles, types.FileStore{
			FileName: msg.fileName,
			FileURL:  msg.fileURL,
			Dir:      msg.filePath,
		})
		m.failedFilesMux.Unlock()

		// Log error to file
		if m.errorLogger != nil {
			m.errorLogger.LogErrorWithDetails(
				errorlog.ErrorTypeDownload,
				fmt.Sprintf("Failed to download file: %s", msg.fileName),
				msg.err,
				map[string]string{
					"file":      msg.fileName,
					"file_url":  msg.fileURL,
					"file_path": msg.filePath,
				},
			)
		}
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "esc" || msg.String() == "ctrl+c" {
			m.cancelled = true
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
		view += lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF8C00")).
			Render(fmt.Sprintf("\nCurrent file: %s\n", m.currentFile))
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
	fileURL  string
	filePath string
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
func DownloadFiles(filesStoreChan <-chan types.FileStore, maxGoroutines int, courses []types.Course, errLogger *errorlog.ErrorLogger) {
	totalFiles := len(filesStoreChan)
	if maxGoroutines == -1 {
		maxGoroutines = totalFiles
	}

	// Convert the courses to a map for easy lookup
	courseIDMap := make(map[string]string)
	for _, course := range courses {
		courseIDMap[course.ID] = course.Name
	}

	m := model{
		totalFiles:  int32(totalFiles),
		courseIDMap: courseIDMap,
		errorLogger: errLogger,
	}

	if m.totalFiles == 0 {
		color.Red("No files to download\n")
		return
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
				if err := downloadFileWithRetry(fileStore, 0); err != nil {
					p.Send(errorMsg{
						fileName: fileStore.FileName,
						fileURL:  fileStore.FileURL,
						filePath: fileStore.Dir,
						err:      err,
					})
				} else {
					p.Send(progressMsg{fileName: fileStore.FileName})
				}
			}(fileStore)
		}
		wg.Wait()

		// Quit the program after all downloads are complete
		p.Send(tea.Quit())
	}()

	finalModel, err := p.Run()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	if finalModel.(model).cancelled {
		os.Exit(0)
	}

	// Retry failed downloads
	finalM := finalModel.(model)
	if len(finalM.failedFiles) > 0 {
		color.Yellow("\nRetrying %d failed downloads...\n", len(finalM.failedFiles))
		retryFailedDownloads(finalM.failedFiles, maxGoroutines, finalM.errorLogger, p)
	}

	color.Green("Download completed successfully \n")
}

// retryFailedDownloads attempts to re-download files that failed on the first attempt
func retryFailedDownloads(failedFiles []types.FileStore, maxGoroutines int, errLogger *errorlog.ErrorLogger, p *tea.Program) {
	if len(failedFiles) == 0 {
		return
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxGoroutines)
	successCount := int32(0)
	failCount := int32(0)

	for _, fileStore := range failedFiles {
		wg.Add(1)
		go func(fileStore types.FileStore) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Try downloading with retry logic
			if err := downloadFileWithRetry(fileStore, 0); err != nil {
				atomic.AddInt32(&failCount, 1)
				if errLogger != nil {
					errLogger.LogErrorWithDetails(
						errorlog.ErrorTypeDownload,
						fmt.Sprintf("Failed to download file after retries: %s", fileStore.FileName),
						err,
						map[string]string{
							"file":      fileStore.FileName,
							"file_url":  fileStore.FileURL,
							"file_path": fileStore.Dir,
							"retries":   "exhausted",
						},
					)
				}
			} else {
				atomic.AddInt32(&successCount, 1)
			}
		}(fileStore)
	}
	wg.Wait()

	if successCount > 0 {
		color.Green("Successfully downloaded %d files on retry\n", successCount)
	}
	if failCount > 0 {
		color.Red("Failed to download %d files even after retries\n", failCount)
	}
}

// downloadFileWithRetry attempts to download a file with exponential backoff retry logic
func downloadFileWithRetry(fileStore types.FileStore, attemptNum int) error {
	err := downloadFile(fileStore)
	if err != nil && attemptNum < maxRetries {
		// Calculate backoff duration (exponential backoff)
		backoffDuration := initialBackoff * time.Duration(1<<uint(attemptNum))
		log.Printf("Download failed for %s (attempt %d/%d), retrying in %v: %v\n",
			fileStore.FileName, attemptNum+1, maxRetries, backoffDuration, err)

		time.Sleep(backoffDuration)
		return downloadFileWithRetry(fileStore, attemptNum+1)
	}
	return err
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
