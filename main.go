package main

import (
	"AGDownloader/parser"
	courseParser "AGDownloader/parser"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	domain     = "aulaglobal.uc3m.es"
	webservice = "/webservice/rest/server.php"
)

type UserInfo struct {
	FullName string `xml:"KEY[name='firstname']>VALUE"`
	UserID   string `xml:"KEY[name='userid']>VALUE"`
}

type Course struct {
	Name string `xml:"SINGLE>KEY[name='fullname']>VALUE"`
	ID   string `xml:"SINGLE>KEY[name='id']>VALUE"`
}

type CourseContent struct {
	Modules []Module `xml:"MULTIPLE>SINGLE>KEY[name='modules']>MULTIPLE>SINGLE"`
}

type Module struct {
	FileURL  string `xml:"KEY[name='fileurl']>VALUE"`
	FileName string `xml:"KEY[name='filename']>VALUE"`
	FileType string `xml:"KEY[name='type']>VALUE"`
}

func getUserInfo(token string) (string, string, error) {
	urlInfo := fmt.Sprintf("https://%s%s?wstoken=%s&wsfunction=core_webservice_get_site_info", domain, webservice, token)
	resp, err := http.Get(urlInfo)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	if strings.Contains(string(body), "invalidtoken") {
		return "", "", fmt.Errorf("Invalid Token")
	}
	fmt.Println("Token is valid")

	// Get the user ID
	start := strings.Index(string(body), "<KEY name=\"username\"><VALUE>")
	end := strings.Index(string(body)[start:], "</VALUE>")
	userID := string(body)[start+28 : start+end]
	// fmt.Println("User ID:", userID)

	// Get the user full name
	// Get the user ID
	start = strings.Index(string(body), "<KEY name=\"fullname\"><VALUE>")
	end = strings.Index(string(body)[start:], "</VALUE>")
	fullname := string(body)[start+28 : start+end]
	// fmt.Println("User ID:", fullname)

	fmt.Printf("Your User ID: %s, %s\n", userID, fullname)
	return userID, fullname, nil
}

func getCourseContent(token, courseID string) ([]Module, error) {
	urlCourse := fmt.Sprintf("https://%s%s?wstoken=%s&wsfunction=core_course_get_contents&courseid=%s", domain, webservice, token, courseID)
	resp, err := http.Get(urlCourse)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var courseContent CourseContent
	if err := xml.Unmarshal(body, &courseContent); err != nil {
		return nil, err
	}

	return courseContent.Modules, nil
}

func saveFiles(token, courseID string, modules []parser.Module, dirPath string) error {
	courseID = strings.ReplaceAll(courseID, "/", "_")
	var path string
	if dirPath == "" {
		path = filepath.Join(os.Getenv("PWD"), "cursos", courseID)
	} else {
		path = filepath.Join(dirPath, courseID)
	}

	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	for i, module := range modules {
		var url string
		if i == 0 {
			url = module.FileURL
		} else {
			url = fmt.Sprintf("%s&token=%s", module.FileURL, token)
		}
		filePath := filepath.Join(path, strings.ReplaceAll(module.FileName, "/", "_"))
		//fmt.Printf("\nDownloading file to: %s\n", filePath)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Couldn't download this file. \n%s\n", err)
			continue
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	color.Blue("Download UC3M Aula Global files from Command Line using 'aulaglobalmovil' Security key")

	token_ptr := flag.String("token", "00", "aulaglobalmobile token used to authenticate")
	flag.Parse()
	token := *token_ptr

	userID, _, err := getUserInfo(token)
	if err != nil {
		color.Red("Error: %s", err)
		return
	}
	fmt.Println("Obtaining courses...")

	courses, err := courseParser.ParseXmlCourses(token, userID)
	if err != nil {
		color.Red("Error: %s", err)
		return
	}

	startTime := time.Now()
	for _, course := range courses {
		fmt.Printf("Course: %s\n", course.Name)
		modules, err := courseParser.GetCourseContent(token, course.ID)
		if err != nil {
			color.Red("Error: %s", err)
			continue
		}
		if err := saveFiles(token, course.Name, modules, "cursosDescargados"); err != nil {
			color.Red("Error: %s", err)
		}
	}
	elapsedTime := time.Since(startTime)
	color.Green("All files downloaded in %s", elapsedTime)
}
