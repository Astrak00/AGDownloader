package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/fatih/color"
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

type Courses struct {
	Courses []Course
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
func getUserInfo(token string) (string, error) {
	url := fmt.Sprintf("https://%s%s?wstoken=%s&wsfunction=core_webservice_get_site_info", domain, webservice, token)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if strings.Contains(string(body), "invalidtoken") {
		return "", fmt.Errorf("invalid token")
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
	return userInfo.UserID, nil
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
		return nil, err
	}

	// Find the course names and ids with two regex and then join them into a Course struct
	courseNames := regexp.MustCompile(`<KEY name="fullname"><VALUE>([^<]+)</VALUE>`)
	courseIDs := regexp.MustCompile(`<KEY name="id"><VALUE>([^<]+)</VALUE>`)
	names := courseNames.FindAllStringSubmatch(string(body), -1)
	ids := courseIDs.FindAllStringSubmatch(string(body), -1)

	var courses Courses
	courses.Courses = make([]Course, 0, len(names))

	for i, name := range names {
		courses.Courses = append(courses.Courses, Course{Name: name[1], ID: ids[i][1]})
	}

	color.Green("Courses found: %d\n", len(courses.Courses))

	return courses.Courses, nil
}

/*
Parses the course for available files and sends them to the channel to be downloaded
*/
func processCourse(course Course, userToken string, dirPath string, errChan chan<- error, filesStoreChan chan<- FileStore) {
	fmt.Printf("Course: %s\n", course.Name)
	files, err := getCourseContent(userToken, course.ID)
	if err != nil {
		fmt.Printf("Error getting course content: %v\n", err)
		errChan <- err
	}
	catalogFiles(course.Name, userToken, files, dirPath, filesStoreChan)
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
	var files []File
	files = make([]File, 0, len(names))

	for i := range names {
		if names[i][1] == "structure" {
			// Insert an empy url to the urls at the position i to fix an empty url error in moodle
			urls = append(urls[:i], append([][]string{{""}}, urls[i:]...)...)
			break
		}
	}
	for i, name := range names {
		if types[i][1] == "file" {
			name_ := name[1]
			type_ := types[i][1]
			url_ := urls[i][1]
			files = append(files, File{FileName: name_, FileURL: url_, FileType: type_})
		}
	}
	// color.Red("Files found: %d\n", len(files))

	return files, nil
}

/*
Formats the files to be downloaded, adding the course name and sends them to the channel
*/
func catalogFiles(courseName string, token string, files []File, dirPath string, filesStoreChan chan<- FileStore) {
	for i, file := range files {
		var url string
		if i == 0 {
			url = file.FileURL
		} else {
			url = file.FileURL + "&token=" + token
		}

		courseName = strings.ReplaceAll(courseName, "/", "-")
		filePath_short := filepath.Join(dirPath, courseName)
		// Join the filePath_short with the filename to create the full path
		filePath := filepath.Join(filePath_short, file.FileName)

		filesStoreChan <- FileStore{FileName: file.FileName, FileURL: url, FileType: file.FileType, Dir: filePath}
	}
}

func downloadFile(fileStore FileStore) error {
	filePath := fileStore.Dir
	fileURL := fileStore.FileURL

	err := downloadFileFromURL(fileURL, filePath)
	if err != nil {
		fmt.Printf("Error downloading the file: %v\n", err)
		return err
	}
	return nil
}

func downloadFileFromURL(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Couldn't download this file. %v\n", err)
		return err
	}
	defer resp.Body.Close()

	// Create the directory path if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		color.Red("Couldn't create directory: %v\n", err)
		return err
	}

	out, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Couldn't create the file: %v\n", err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Couldn't copy the file: %v\n", err)
		return err
	}

	//fmt.Printf("Downloaded file: %s\n", filePath)
	return nil
}

func main() {
	blue := color.New(color.FgBlue, color.Bold, color.Underline).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	fmt.Println("Download UC3M Aula Global files from Command Line using 'aulaglobalmovil' Security key")

	// Parse the flags with flag package
	// -l for language
	// -t for token
	// -d for directory

	var nFlag = flag.Int("l", 1, "Choose your language / Escoga su Idioma: 1: Español, 2:English.")
	var tFlag = flag.String("t", "", "Introduce your Aula Global user security token 'aulaglobalmovil'")
	var dFlag = flag.String("d", "courses", "Introduce the directory where you want to save the files")
	var cFlag = flag.Int("c", 7, "Introduce the cores to use in the download")

	flag.Parse()

	var language int
	if *nFlag != 1 && *nFlag != 2 {
		for language != 1 && language != 2 {
			fmt.Println("Choose your language / Escoga su Idioma:")
			fmt.Print("1: Español, 2:English. :")
			_, err := fmt.Scanf("%d", &language)
			if err != nil {
				fmt.Println(red("Wrong input value / Valor introducido erroneo. Intentelo de nuevo"))
			}
			fmt.Println()
		}
	} else {
		language = *nFlag
	}

	var userToken string = *tFlag
	var dirPath string = *dFlag
	var done bool = false
	if userToken != "" {
		done = true
	}
	if dirPath == "" {
		dirPath = "courses"
	}

	for !done {
		if language == 1 {
			fmt.Println("Introduzca el token de seguridad de su usuario de Aula Global 'aulaglobalmovil':")
			fmt.Printf("Para ver las instrucciones para generar el token ve a: %s\n", blue("https://github.com/Josersanvil/AulaGlobal-CoursesFiles#para-conseguir-el-token-de-seguridad"))
			fmt.Print("Introduzca su token se seguridad: ")
		} else {
			fmt.Println("Introduce your Aula Global user security token 'aulaglobalmovil':")
			fmt.Printf("To see instructions on how to generate the token go to: %s\n", blue("https://github.com/Josersanvil/AulaGlobal-CoursesFiles#get-your-token"))
			fmt.Print("Introduce your security token: ")
		}

		fmt.Scanf("%s", &userToken)

		if userToken != "" && regexp.MustCompile(`[a-zA-Z0-9]{20,}`).MatchString(userToken) && len(userToken) > 20 {
			done = true
		} else {
			if language == 1 {
				fmt.Println(red("El token introducido no parece estar correcto. Intentalo de nuevo."))
			} else {
				fmt.Println(red("The given token does not seem to be right. Please try again."))
			}
		}
		fmt.Println()
	}

	if language == 1 {
		fmt.Println("Descargando los archivos a la carpeta", blue(dirPath))
	} else {
		fmt.Println("Downloading files to the folder", blue(dirPath))
	}

	userID, err := getUserInfo(userToken)
	if err != nil {
		if language == 1 {
			fmt.Println(red("Error: Token invalido, el token podria haber expirado o es erroneo. Chequea que esta escrito correctamente o generara uno nuevo en 'aulaglobal.uc3m.es' > perfil."))
		} else {
			fmt.Println(red("Error: Invalid Token, the token may have expired or has a typo. Check if it's written correctly or generate a new one in 'aulaglobal.uc3m.es' > profile."))
		}
		os.Exit(1)
	}

	courses, err := getCourses(userToken, userID)
	if err != nil {
		fmt.Printf("Error getting courses: %v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(courses))
	filesStoreChan := make(chan FileStore, 250)

	for _, course := range courses {
		wg.Add(1)
		go func(course Course) {
			defer wg.Done()
			var (
				userToken string = userToken
				dirPath   string = dirPath
			)
			fmt.Printf("Course: %s\n", course.Name)
			files, err := getCourseContent(userToken, course.ID)
			if err != nil {
				fmt.Printf("Error getting course content: %v\n", err)
				chan<- error(errChan) <- err
			}
			catalogFiles(course.Name, userToken, files, dirPath, chan<- FileStore(filesStoreChan))
		}(course)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			fmt.Println("Error:", err)
		}
	}

	close(filesStoreChan)
	color.Red("Found %d items to download\n", len(filesStoreChan))

	// Create the directory to store the files
	err = os.MkdirAll(dirPath, 0755)
	if err != nil {
		fmt.Printf("Error creating the directory: %v\n", err)
		os.Exit(1)
	}

	// Create a subdirectory for each course
	for _, course := range courses {
		courseName := strings.ReplaceAll(course.Name, "/", "-")
		if len(courseName) > 55 {
			courseName = courseName[:55]
		}

		courseDir := filepath.Join(dirPath, courseName)
		err = os.MkdirAll(courseDir, 0755)
		if err != nil {
			fmt.Printf("Error creating the directory: %v\n", err)
			os.Exit(1)
		}
	}
	color.Green("Downloading files...\n")

	// Download the files from the channel by goroutines
	maxGoroutines := *cFlag
	var wg2 sync.WaitGroup
	errChan2 := make(chan error, len(courses)*20)
	semaphore := make(chan struct{}, maxGoroutines)

	for fileStore := range filesStoreChan {
		semaphore <- struct{}{} // Acquire a slot
		wg2.Add(1)
		go func(fileStore FileStore) {
			defer wg2.Done()
			defer func() { <-semaphore }() // Release the slot

			err := downloadFile(fileStore)
			if err != nil {
				errChan2 <- err
			}
		}(fileStore)
	}

	wg2.Wait()
	close(errChan2)
}
