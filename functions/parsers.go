package functions

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	domain     = "aulaglobal.uc3m.es"
	webservice = "/webservice/rest/server.php"
)

type Response struct {
	XMLName  xml.Name `xml:"RESPONSE"`
	Multiple Multiple `xml:"MULTIPLE"`
}

type Multiple struct {
	Singles []Single `xml:"SINGLE"`
}

type Single struct {
	Keys []Key `xml:"KEY"`
}

type Key struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"VALUE"`
}

type Course struct {
	ID   string
	Name string
}

type Module struct {
	FileURL  string
	FileName string
}

func ParseXmlCourses(token, userID string) ([]Course, error) {
	urlCourses := fmt.Sprintf("https://%s%s?wstoken=%s&wsfunction=core_enrol_get_users_courses&userid=%s", domain, webservice, token, userID)
	resp, err := http.Get(urlCourses)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	byteValue, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	err = xml.Unmarshal(byteValue, &response)
	if err != nil {
		log.Fatal(err)
	}

	courses := []Course{}

	// Iterate through the parsed XML and extract 'id' and 'shortname'
	for _, single := range response.Multiple.Singles {
		var id, name string
		for _, key := range single.Keys {
			if key.Name == "id" {
				id = key.Value
			}
			if key.Name == "fullname" {
				name = key.Value
			}
		}
		if id != "" && name != "" {
			courses = append(courses, Course{ID: id, Name: name})
		}
	}
	return courses, nil

}

func GetCourseContent(token, courseID string) ([]Module, error) {
	urlCourse := fmt.Sprintf("https://%s%s?wstoken=%s&wsfunction=core_course_get_contents&courseid=%s", domain, webservice, token, courseID)
	resp, err := http.Get(urlCourse)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	byteValue, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response Response
	err = xml.Unmarshal(byteValue, &response)
	if err != nil {
		log.Fatal(err)
	}

	modules := []Module{}

	// Cada single es una sección/tema de AulaGlobal. (En inglés es topic)

	// Iterate through the parsed XML and extract 'fileurl', 'filename' and 'type'
	// Search for the KEY with name "url" and "filename" manually and then append it to the modules slice
	xmlContent := string(byteValue)

	start := strings.Index(xmlContent, "fileurl\"><VALUE>")
	if start == -1 {
		return nil, fmt.Errorf("no modules found")
	}
	end := 0
	start_new := 0

	for {
		// Get the user ID
		start_new = strings.Index(xmlContent[start:], "<KEY name=\"fileurl\"><VALUE>")
		end = strings.Index(xmlContent[start+start_new:], "</VALUE>")
		url := xmlContent[start+start_new+27 : start+start_new+end]
		// fmt.Print("URL:", url)
		start = start + start_new + end

		// fmt.Println("User ID:", userID)
		// xmlContent = xmlContent[start:]

		// Get the file name

		start_new = strings.Index(xmlContent[start:], "<KEY name=\"filename\"><VALUE>")
		if start_new <= 0 {
			break
		}
		end = strings.Index(xmlContent[start+start_new:], "</VALUE>")
		filename := xmlContent[start+start_new+28 : start+start_new+end]
		// fmt.Println(" -- Filename:", filename)
		start = start + start_new + end

		modules = append(modules, Module{FileURL: url, FileName: filename})

		if start_new <= 0 {
			break
		}
	}

	return modules, nil
}

func GetUserInfo(token string) (string, string, error) {
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