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
	color.Yellow("Fetching courses from AulaGlobal...\n")

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
		if name_def != "Secretaría EPS" && !containsInvalidNames(name_def) {
			courses = append(courses, types.Course{Name: name_def, ID: strconv.Itoa(course.ID)})
		}
	}

	defer color.Green("Number of courses found: %d\n", len(courses))
	return courses, nil
}

func containsInvalidNames(name string) bool {
	invalidCourseNames := []string{"Convenio", "Delegación", "Secretaría", "Student Room", "Sala de Estudiantes", "Bachelor"}

	for _, invalidName := range invalidCourseNames {
		if strings.Contains(name, invalidName) {
			return true
		}
	}
	return false

}

// Get the names of the courses in Spanish and English
// This function localizes the names of the courses in Spanish and English
// Separating the names by -1C, -2C, -1S, -2S, Bachelor, Student, Convenio-Bilateral
func getCoursesNamesLanguages(name string) (string, string) {
	// Define the first group of separators with priority.
	firstGroup := []string{"-1C", "-2C", "-1S", "-2S"}

	// Iterate over the first group to find the earliest separator.
	for _, sep := range firstGroup {
		if idx := strings.Index(name, sep); idx > 0 { // idx > 0 ensures the separator is not at the start
			return name[:idx+len(sep)], name[idx+len(sep):]
		}
	}

	// Define the second group of separators.
	secondGroup := []string{"Bachelor", "Student", "Convenio-Bilateral s"}

	// Iterate over the second group to find the earliest separator.
	for _, sep := range secondGroup {
		if idx := strings.Index(name, sep); idx != -1 { // idx != -1 means the separator exists
			return name[:idx], name[idx:]
		}
	}

	// If no separators are found, return the original name twice.
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
