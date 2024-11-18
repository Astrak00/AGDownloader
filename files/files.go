package files

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	types "github.com/Astrak00/AGDownloader/types"
)

// Parses the course for available files and sends them to the channel to be downloaded
func processCourse(course types.Course, userToken string, dirPath string, errChan chan<- error, filesStoreChan chan<- types.FileStore) {
	files, err := getCourseContent(userToken, course.ID)
	if err != nil {
		errChan <- fmt.Errorf("error getting course content: %v", err)
	}
	if len(files) > 0 {
		catalogFiles(course.Name, userToken, files, dirPath, filesStoreChan)
	}
}

// Parses the course and returns the files of type "file"
// Fetches the course content from the moodle API
// Scrapes the file names, urls and types with regex
func getCourseContent(token, courseID string) ([]types.File, error) {
	url := fmt.Sprintf("https://%s%s?wstoken=%s&wsfunction=core_course_get_contents&courseid=%s", types.Domain, types.Webservice, token, courseID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing body")
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing body")
		}
	}(resp.Body)

	// Find the file names, urls and types with regex
	fileNames := regexp.MustCompile(`<KEY name="filename"><VALUE>([^<]+)</VALUE>`)
	fileURLs := regexp.MustCompile(`<KEY name="fileurl"><VALUE>([^<]+)</VALUE>`)
	fileTypes := regexp.MustCompile(`<KEY name="type"><VALUE>([^<]+)</VALUE>`)

	names := fileNames.FindAllStringSubmatch(string(body), -1)
	// If no files are found, return nil and prevent further execution and useless processing
	if len(names) == 0 {
		// color.Red("No files found\n")
		return nil, nil
	}
	urls := fileURLs.FindAllStringSubmatch(string(body), -1)
	fileType := fileTypes.FindAllStringSubmatch(string(body), -1)

	// Insert an empy url to the urls at the position i to fix an empty url error in moodle
	for i := range names {
		if names[i][1] == "structure" {
			urls = append(urls[:i], append([][]string{{""}}, urls[i:]...)...)
			break
		}
	}

	// Join the names and urls into a File struct
	files := make([]types.File, 0, len(names))
	for i, name := range names {
		if fileType[i][1] == "file" {
			files = append(files, types.File{
				FileName: name[1],
				FileURL:  urls[i][1],
				FileType: fileType[i][1],
			})
		}
	}
	return files, nil
}

// Formats the files to be downloaded, adding the course name and sends them to the channel
func catalogFiles(courseName string, token string, files []types.File, dirPath string, filesStoreChan chan<- types.FileStore) {
	for _, file := range files {
		url := file.FileURL + "&token=" + token

		// Replace the "/" in the course name to avoid creating subdirectories
		courseName = strings.ReplaceAll(courseName, "/", "-")
		filePath := filepath.Join(dirPath, courseName, file.FileName)

		// Send the file to the channel
		filesStoreChan <- types.FileStore{FileName: file.FileName, FileURL: url, FileType: file.FileType, Dir: filePath}
	}
}

// ListAllResourcess Creates a list of all the resources to download
func ListAllResourcess(downloadAll bool, courses []types.Course, userToken string, dirPath string, errChan chan error, filesStoreChan chan types.FileStore, coursesList []string) {
	var wg sync.WaitGroup
	if downloadAll {
		for _, courseItem := range courses {
			wg.Add(1)
			go func(courseItem types.Course) {
				defer wg.Done()
				// Passing chan <- types.FileStore(filesStoreChan) as a parameter to the function makes the chanel
				// to be a parameter of the function, so it can be used inside the function and a send-only channel
				processCourse(courseItem, userToken, dirPath, chan<- error(errChan), chan<- types.FileStore(filesStoreChan))
			}(courseItem)
		}
	} else {
		for _, course := range courses {
			for _, courseSearch := range coursesList {
				if courseSearch == course.ID || strings.Contains(strings.ToLower(course.Name), strings.ToLower(courseSearch)) {
					wg.Add(1)
					go func(course types.Course) {
						defer wg.Done()
						processCourse(course, userToken, dirPath, chan<- error(errChan), chan<- types.FileStore(filesStoreChan))
					}(course)
				}
			}
		}
	}
	wg.Wait()
}
