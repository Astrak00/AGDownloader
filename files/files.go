package files

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	types "github.com/Astrak00/AGDownloader/types"
)

// ListAllResources Creates a list of all the resources to download
func ListAllResources(courses []types.Course, userToken string, dirPath string, errChan chan error, filesStoreChan chan types.FileStore) {
	var wg sync.WaitGroup
	for _, courseItem := range courses {
		wg.Add(1)
		go func(courseItem types.Course) {
			defer wg.Done()
			// Passing chan <- types.FileStore(filesStoreChan) as a parameter to the function makes the chanel
			// to be a parameter of the function, so it can be used inside the function and a send-only channel
			processCourse(courseItem, userToken, dirPath, chan<- error(errChan), chan<- types.FileStore(filesStoreChan))
		}(courseItem)
	}

	wg.Wait()
}

// Parses the course for available files and sends them to the channel to be downloaded
func processCourse(course types.Course, userToken string, dirPath string, errChan chan<- error, filesStoreChan chan<- types.FileStore) {
	files, err := getCourseContent(userToken, course.ID)
	if err != nil {
		errChan <- fmt.Errorf("error getting course content: %v", err)
	}
	if len(files) > 0 {
		// Replace the "/" in the course name to avoid creating subdirectories
		courseName := strings.ReplaceAll(course.Name, "/", "-")
		catalogFiles(courseName, userToken, files, dirPath, filesStoreChan)
	}
}

// Parses the course and returns the files of type "file"
// Fetches the course content from the moodle API
// Scrapes the file names, urls and types with regex
func getCourseContent(token, courseID string) ([]types.File, error) {
	url := fmt.Sprintf("https://%s%s?wstoken=%s&wsfunction=core_course_get_contents&moodlewsrestformat=json&courseid=%s", types.Domain, types.Webservice, token, courseID)

	// Get the json from the URL
	jsonData := types.GetJson(url)

	// Parse the json
	var courseParsed types.WebCourse
	err := json.Unmarshal(jsonData, &courseParsed)
	if err != nil {
		log.Fatal(err)
	}

	// Get the names, urls and types of the files
	filesPresentInCourse := make([]types.File, 0)
	for _, course := range courseParsed {
		if len(course.Modules) == 0 {
			continue
		}

		// get section name
		sectionName := course.Name

		if sectionName == "General" {
			// save it without section
			sectionName = ""
		}

		if strings.HasPrefix(sectionName, "Topic ") || strings.HasPrefix(sectionName, "Tema ") { // TODO: use one or the other depending on the current language
			// the section has a generic name, search the name in the summary
			sectionName = removeTags(course.Summary)
		}

		for _, module := range course.Modules {
			for _, content := range module.Contents {

				switch content.Type {
				case "file":
					if runtime.GOOS == "windows" {
						// Windows shits itself if the sectionName has trailing whitespace or `:`, don't ask me why (see #13)
						sectionName = strings.Replace(strings.TrimSpace(sectionName), ":", "", -1)
					}
					filesPresentInCourse = append(filesPresentInCourse, types.File{
						FileName: filepath.Join(sectionName, content.Filename),
						FileURL:  content.Fileurl,
					})
				default:
					continue
				}
			}
		}
	}

	return filesPresentInCourse, nil
}

func removeTags(s string) string {
	// Remove tags from a string using a more efficient approach

	s = strings.ReplaceAll(s, "&nbsp;", " ") // remove &nsbp

	var result []rune
	inTag := false
	for _, r := range s {
		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case r == '\n' || r == '\t' || r == '\r':
			continue
		default:
			if !inTag {
				result = append(result, r)
			}
		}
	}
	return strings.Trim(string(result), " ")
}

// Formats the files to be downloaded, adding the course name and sends them to the channel
func catalogFiles(courseName string, token string, files []types.File, dirPath string, filesStoreChan chan<- types.FileStore) {
	for _, file := range files {
		url := file.FileURL + "&token=" + token
		filePath := filepath.Join(dirPath, courseName, file.FileName)

		// Send the file to the channel
		filesStoreChan <- types.FileStore{FileName: file.FileName, FileURL: url, Dir: filePath}
	}
}
