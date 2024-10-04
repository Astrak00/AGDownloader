package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	types "github.com/Astrak00/AGDownloader/types"
	"github.com/fatih/color"

	"github.com/schollz/progressbar/v3"
)

// Download the files in the channel, with a maximum of goroutines and a language
// Indicates with a progress bar the download of the files
func DownloadFiles(filesStoreChan <-chan types.FileStore, maxGoroutines int, language int) {
	if language == 1 {
		color.Red("Se han encontrado %d archivos para descargar\n", len(filesStoreChan))
	} else {
		color.Red("Found %d items to download\n", len(filesStoreChan))
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxGoroutines)
	totalFiles := len(filesStoreChan)

	// Create an atomic counter for completed files
	var completedFiles int32

	// Create a progress bar
	bar := progressbar.NewOptions(totalFiles,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	for fileStore := range filesStoreChan {
		wg.Add(1)
		go func(fileStore types.FileStore) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			if err := downloadFileWithProgress(fileStore, bar, &completedFiles); err != nil {
				fmt.Printf("Error downloading file %s: %v\n", fileStore.FileName, err)
			}
		}(fileStore)
	}
	wg.Wait()
}

// Download the file with a progress bar, meaning when the file is downloaded, the progress bar will increase
// As the counter is passed by reference, it will increase the number of completed files.
func downloadFileWithProgress(fileStore types.FileStore, bar *progressbar.ProgressBar, completedFiles *int32) error {
	resp, err := http.Get(fileStore.FileURL)
	if err != nil {
		return fmt.Errorf("error downloading the file: %v", err)
	}
	defer resp.Body.Close()

	dir := filepath.Dir(fileStore.Dir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating the directory: %v", err)
	}

	out, err := os.Create(fileStore.Dir)
	if err != nil {
		return fmt.Errorf("error creating the file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error copying the file: %v", err)
	}

	atomic.AddInt32(completedFiles, 1)
	bar.Add(1)
	return nil
}
