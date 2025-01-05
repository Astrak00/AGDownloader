package files

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
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
	for i := 0; i < len(courseParsed); i++ {
		if len(courseParsed[i].Modules) != 0 {
			sectionName := removeTags(courseParsed[i].Summary)

			for j := 0; j < len(courseParsed[i].Modules); j++ {
				if len(courseParsed[i].Modules[j].Contents) != 0 {
					for k := 0; k < len(courseParsed[i].Modules[j].Contents); k++ {
						if courseParsed[i].Modules[j].Contents[k].Type == "file" {
							filesPresentInCourse = append(filesPresentInCourse, types.File{
								FileName: filepath.Join(sectionName, courseParsed[i].Modules[j].Contents[k].Filename),
								FileURL:  courseParsed[i].Modules[j].Contents[k].Fileurl,
							})
						}
					}
				}
			}
		}
	}

	return filesPresentInCourse, nil
}

func removeTags(s string) string {
	// Remove tags from a string using a more efficient approach
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

// ListAllResourcess Creates a list of all the resources to download
func ListAllResourcess(courses []types.Course, userToken string, dirPath string, errChan chan error, filesStoreChan chan types.FileStore) {
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
