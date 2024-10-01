package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

const (
	domain     = "aulaglobal.uc3m.es"
	webservice = "/webservice/rest/server.php"
	service    = "aulaglobal_mobile"
)

type UserInfo struct {
	FullName string
	UserID   string
}

type Course struct {
	Name string
	ID   string
}

type File struct {
	FileName string
	FileURL  string
	FileType string
}

type FileStore struct {
	FileName string
	FileURL  string
	FileType string
	Dir      string
}

/*
Gets the userID necessary to get the courses
*/
func getUserInfo(token string) (UserInfo, error) {
	url := fmt.Sprintf("https://%s%s?wstoken=%s&wsfunction=core_webservice_get_site_info", domain, webservice, token)

	resp, err := http.Get(url)
	if err != nil {
		return UserInfo{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserInfo{}, err
	}

	if strings.Contains(string(body), "invalidtoken") {
		return UserInfo{}, fmt.Errorf("invalid token")
	}

	var userInfo UserInfo

	// Find the fullname key and value
	fullName := regexp.MustCompile(`<KEY name="fullname"><VALUE>([^<]+)</VALUE>`)
	maches := fullName.FindStringSubmatch(string(body))
	if len(maches) > 1 {
		userInfo.FullName = maches[1]
	} else {
		color.Red("Fullname not found\n")
	}

	// Find the userid key and value
	userID := regexp.MustCompile(`<KEY name="userid"><VALUE>([^<]+)</VALUE>`)
	maches = userID.FindStringSubmatch(string(body))
	if len(maches) > 1 {
		userInfo.UserID = maches[1]
	} else {
		color.Red("UserID not found\n")
	}

	//color.Blue("Your User ID: %s, %s\n", userInfo.UserID, userInfo.FullName)
	return userInfo, nil
}

/*
Gets the courses, both name and ID, of a given userID
*/
func getCourses(token, userID string) ([]Course, error) {
	url := fmt.Sprintf("https://%s%s?wstoken=%s&wsfunction=core_enrol_get_users_courses&userid=%s", domain, webservice, token, userID)
	color.Yellow("Getting courses...\n")

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("mismatch in course names and IDs")
	}

	// Find the course names and ids with two regex and then join them into a Course struct
	courseNames := regexp.MustCompile(`<KEY name="fullname"><VALUE>([^<]+)</VALUE>`)
	courseIDs := regexp.MustCompile(`<KEY name="id"><VALUE>([^<]+)</VALUE>`)
	names := courseNames.FindAllStringSubmatch(string(body), -1)
	ids := courseIDs.FindAllStringSubmatch(string(body), -1)

	courses := make([]Course, 0, len(names))
	for i, name := range names {
		courses = append(courses, Course{Name: name[1], ID: ids[i][1]})
	}
	color.Green("Courses found: %d\n", len(courses))

	return courses, nil
}

/*
Parses the course for available files and sends them to the channel to be downloaded
*/
func processCourse(course Course, userToken string, dirPath string, errChan chan<- error, filesStoreChan chan<- FileStore) {
	//fmt.Printf("Course: %s\n", course.Name)
	files, err := getCourseContent(userToken, course.ID)
	if err != nil {
		errChan <- fmt.Errorf("error getting course content: %v", err)
	}
	if len(files) > 0 {
		catalogFiles(course.Name, userToken, files, dirPath, filesStoreChan)
	}
}

/*
Parses the course and returns the files of type "file"
*/
func getCourseContent(token, courseID string) ([]File, error) {
	url := fmt.Sprintf("https://%s%s?wstoken=%s&wsfunction=core_course_get_contents&courseid=%s", domain, webservice, token, courseID)
	//color.Cyan("Getting course content...\n")
	//color.Cyan("URL: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
	types := fileTypes.FindAllStringSubmatch(string(body), -1)

	// Join the names and urls into a File struct
	files := make([]File, 0, len(names))

	for i := range names {
		if names[i][1] == "structure" {
			// Insert an empy url to the urls at the position i to fix an empty url error in moodle
			urls = append(urls[:i], append([][]string{{""}}, urls[i:]...)...)
			break
		}
	}
	for i, name := range names {
		if types[i][1] == "file" {
			files = append(files, File{
				FileName: name[1],
				FileURL:  urls[i][1],
				FileType: types[i][1],
			})
		}
	}
	// color.Red("Files found: %d\n", len(files))

	return files, nil
}

/*
Formats the files to be downloaded, adding the course name and sends them to the channel
*/
func catalogFiles(courseName string, token string, files []File, dirPath string, filesStoreChan chan<- FileStore) {
	for _, file := range files {
		url := file.FileURL + "&token=" + token

		// Replace the "/" in the course name to avoid creating subdirectories
		courseName = strings.ReplaceAll(courseName, "/", "-")
		filePath := filepath.Join(dirPath, courseName, file.FileName)

		// Send the file to the channel
		filesStoreChan <- FileStore{FileName: file.FileName, FileURL: url, FileType: file.FileType, Dir: filePath}
	}
}

func downloadFile(fileStore FileStore) error {
	filePath := fileStore.Dir
	fileURL := fileStore.FileURL

	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("error downloading the file: %v", err)
	}
	defer resp.Body.Close()

	// Create the directory path if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating the directory: %v", err)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating the file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error copying the file: %v", err)
	}
	return nil
}

func main() {
	language, userToken, dirPath, maxGoroutines := parseFlags()

	if language == 1 {
		color.Cyan("Programa creado por Astrak00: github.com/Astrak00/AGDownloader/ para descargar archivos de Aula Global en la UC3M\n")
		fmt.Println("Descargando los archivos a la carpeta", color.BlueString(dirPath))
	} else {
		color.Cyan("Program created by Astrak00: github.com/Astrak00/AGDownloader/ to download files from Aula Global at UC3M\n")
		fmt.Println("Downloading files to the folder", color.BlueString(dirPath))
	}

	user, err := getUserInfo(userToken)
	if err != nil {
		log.Fatalf("Error: Invalid Token: %v", err)
	}

	courses, err := getCourses(userToken, user.UserID)
	if err != nil {
		log.Fatalf("Error getting courses: %v\n", err)
	}

	filesStoreChan := make(chan FileStore, len(courses)*20)
	errChan := make(chan error, len(courses))

	var wg sync.WaitGroup
	for _, course := range courses {
		wg.Add(1)
		go func(course Course) {
			defer wg.Done()
			processCourse(course, userToken, dirPath, chan<- error(errChan), chan<- FileStore(filesStoreChan))
		}(course)
	}

	wg.Wait()
	close(errChan)
	close(filesStoreChan)

	for err := range errChan {
		if err != nil {
			fmt.Println("Error:", err)
		}
	}

	if language == 1 {
		color.Red("Se han encontrado %d archivos para descargar\n", len(filesStoreChan))
	} else {
		color.Red("Found %d items to download\n", len(filesStoreChan))
	}

	downloadFiles(filesStoreChan, maxGoroutines)

	if language == 1 {
		color.Green("Descarga completada\n")
	} else {
		color.Green("Download completed\n")
	}

}

func parseFlags() (int, string, string, int) {
	language := flag.Int("l", 0, "Choose your language: 1: Español, 2:English")
	token := flag.String("t", "", "Introduce your Aula Global user security token 'aulaglobalmovil'")
	dir := flag.String("d", "courses", "Introduce the directory where you want to save the files")
	cores := flag.Int("c", runtime.NumCPU()-1, "Introduce the cores to use in the download")
	flag.Parse()

	if *language == 0 {
		fmt.Println("Introduce tu idioma: 1: Español, 2:English")
		fmt.Scanf("%d", language)
	}

	if *token == "" {
		*token = promptForToken(*language)
	}

	return *language, *token, *dir, *cores
}

func promptForToken(language int) string {
	var token string
	for {
		if language == 1 {
			fmt.Println("Introduzca el token de seguridad de su usuario de Aula Global 'aulaglobalmovil':")
		} else {
			fmt.Println("Introduce your Aula Global user security token 'aulaglobalmovil':")
		}
		fmt.Scanf("%s", &token)

		if token != "" && regexp.MustCompile(`[a-zA-Z0-9]{20,}`).MatchString(token) && len(token) > 20 {
			return token
		}

		if language == 1 {
			color.Red("El token introducido no parece estar correcto. Inténtelo de nuevo.")
		} else {
			color.Red("The given token does not seem to be right. Please try again.")
		}
	}
}

func downloadFiles(filesStoreChan <-chan FileStore, maxGoroutines int) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxGoroutines)
	totalFiles := len(filesStoreChan)

	// Create an atomic counter for completed files
	var completedFiles int32

	bar := progressbar.NewOptions(totalFiles,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(20),
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
		go func(fileStore FileStore) {
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

func downloadFileWithProgress(fileStore FileStore, bar *progressbar.ProgressBar, completedFiles *int32) error {
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
