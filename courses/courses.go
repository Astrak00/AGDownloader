package courses

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"

	types "github.com/Astrak00/AGDownloader/types"
)

// GetCourses Gets the courses, the localized name and ID, given a userID
func GetCourses(token string, userID string, language int) ([]types.Course, error) {
	if language == 1 {
		color.Yellow("Obteniendo cursos de AulaGlobal...\n")
	} else {
		color.Yellow("Getting courses from AulaGlobal...\n")
	}

	url := fmt.Sprintf("https://%s%s?wstoken=%s&wsfunction=core_enrol_get_users_courses&userid=%s&moodlewsrestformat=json", types.Domain, types.Webservice, token, userID)

	jsonData := types.GetJson(url)

	// Parse the json
	var userParsed types.WebUser
	err := json.Unmarshal(jsonData, &userParsed)
	if err != nil {
		log.Fatal(err)
	}

	// Get the names and IDs of the courses
	courses := make([]types.Course, 0, len(userParsed))
	name_def := ""
	for _, course := range userParsed {
		// Localize the course name
		nameEs, nameEN := getCoursesNamesLanguages(course.Fullname)
		if language == 1 {
			name_def = nameEs
		} else {
			name_def = nameEN
		}
		if name_def != "Secretaría EPS" &&
			!strings.Contains(name_def, "Convenio") &&
			!strings.Contains(name_def, "Delegación") {
			courses = append(courses, types.Course{Name: name_def, ID: strconv.Itoa(course.ID)})
		}
	}

	defer func() {
		if language == 1 {
			color.Green("Cursos encontrados: %d\n", len(courses))
		} else {
			color.Green("Courses found: %d\n", len(courses))
		}
	}()

	return courses, nil
}

// Get the names of the courses in Spanish and English
// This function localizes the names of the courses in Spanish and English
// Separating the names by -1C or -2C
func getCoursesNamesLanguages(name string) (string, string) {
	// Find where the names are separated, by -1C or -2C and return the names in Spanish and English
	idx := 0
	if strings.Contains(name, "-1C") {
		idx = strings.Index(name, "-1C")
	} else if strings.Contains(name, "-2C") {
		idx = strings.Index(name, "-2C")
	} else if strings.Contains(name, "-1S") {
		idx = strings.Index(name, "-1S")
	} else if strings.Contains(name, "-2S") {
		idx = strings.Index(name, "-2S")
	}
	if idx != 0 {
		return name[:idx+3], name[idx+3:]
	}

	if strings.Contains(name, "Bachelor") {
		idx = strings.Index(name, "Bachelor")
		return name[:idx], name[idx:]
	} else if strings.Contains(name, "Student") {
		idx = strings.Index(name, "Student")
		return name[:idx], name[idx:]
	} else if strings.Contains(name, "Convenio-Bilateral s") {
		idx = strings.Index(name, "Convenio-Bilateral s")
		return name[:idx], name[idx:]
	}
	return name, name
}

// SelectCourses prompts the user to select the courses to download
func SelectCourses(language int, coursesList []string, courses []types.Course) []types.Course {
	// Cehck if the user wants to download all the courses and return the courses
	if len(coursesList) != 0 && coursesList[0] == "all" {
		return courses
	} else if len(coursesList) == 0 {
		// In case the user does not want to download all the courses, show a list of checkboxes with the courses
		// to allow the user to select them interactively
		prompt := "Select the courses you want to download\n"
		if language == 1 {
			prompt = "Selecciona los cursos que quieres descargar\n"
		}
		listCoursesList := getCoursesNameByLanguage(courses)
		coursesList = checkboxesCourses(prompt, listCoursesList)
	}

	coursesToDownload := make([]types.Course, 0, len(coursesList))
	for _, course := range coursesList {
		for _, c := range courses {
			if course == c.Name {
				coursesToDownload = append(coursesToDownload, c)
			}
		}
	}
	return coursesToDownload
}

// Map the courses to obtain a []string with the names of the courses
func getCoursesNameByLanguage(courses []types.Course) []string {
	coursesList := make([]string, len(courses))
	for i, course := range courses {
		coursesList[i] = course.Name
	}
	return coursesList
}

// Show in the terminal a list of checkboxes with the courses to download
func checkboxesCourses(label string, opts []string) []string {
	res := []string{}
	prompt := &survey.MultiSelect{
		Message:  label,
		Options:  opts,
		PageSize: 6,
	}
	err := survey.AskOne(prompt, &res, survey.WithKeepFilter(true))
	if err != nil {
		return nil
	}

	return res
}
