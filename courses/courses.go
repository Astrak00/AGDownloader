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

// GetCourses obtains the courses, the localized name and ID, given a userID
// Returns a slice of courses
func GetCourses(token string, userID string, language int) (types.Courses, error) {
	fmt.Println("Fetching courses from AulaGlobal...")

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
	var course_name string = ""
	for _, course := range userParsed {
		// Localize the course name
		course_name = extractCourseNameByLanguage(course.Fullname, language)
		if !containsInvalidNames(course_name) {
			courses = append(courses, types.Course{Name: course_name, ID: strconv.Itoa(course.ID)})
		}
	}

	defer color.Green("Number of courses found: %d\n", len(courses))
	return courses, nil
}

// Check if the name of the course contains invalid names that should not be downloaded
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
func extractCourseNameByLanguage(name string, lang int) string {
	// Define the first group of separators with priority.
	firstGroup := []string{"-1C", "-2C", "-1S", "-2S"}

	// Iterate over the first group to find the earliest separator.
	for _, sep := range firstGroup {
		if idx := strings.Index(name, sep); idx > 0 { // idx > 0 ensures the separator is not at the start
			if lang == 1 {
				return name[:idx+len(sep)]
			} else {
				return name[idx+len(sep):]
			}
		}
	}

	// Define the second group of separators.
	secondGroup := []string{"Bachelor", "Student", "Convenio-Bilateral s"}

	// Iterate over the second group to find the earliest separator.
	for _, sep := range secondGroup {
		if idx := strings.Index(name, sep); idx != -1 { // idx != -1 means the separator exists
			if lang == 1 {
				return name[:idx]
			} else {
				return name[idx:]
			}
		}
	}

	// If no separators are found, return the original name twice.
	return name
}

// SelectCourses prompts the user to select the courses to download
func SelectCoursesInteractive(language int, selectedCourses []string, courses types.Courses) []types.Course {
	// Cehck if the user wants to download all the courses and return the courses
	if len(selectedCourses) != 0 && selectedCourses[0] == "all" {
		return courses
	} else if len(selectedCourses) == 0 {
		// In case the user does not want to download all the courses, show a list of checkboxes with the courses
		// to allow the user to select them interactively
		prompt := "Select the courses you want to download\n"

		coursesName := courses.GetCoursesName()
		selectedCourses = checkboxesCourses(prompt, coursesName)
	}

	coursesToDownload := make([]types.Course, 0, len(selectedCourses))
	// Create a map for O(1) lookups
	courseMap := make(map[string]types.Course)
	for _, c := range courses {
		courseMap[c.Name] = c
	}

	// Single loop through selected courses
	for _, courseName := range selectedCourses {
		if course, exists := courseMap[courseName]; exists {
			coursesToDownload = append(coursesToDownload, course)
		}
	}
	return coursesToDownload
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
