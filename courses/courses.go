package courses

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"

	types "github.com/Astrak00/AGDownloader/types"
)

// Gets the courses, the localized name and ID, given a userID
func GetCourses(token string, userID string, language int) ([]types.Course, error) {
	if language == 1 {
		color.Yellow("Obteniendo cursos...\n")
	} else {
		color.Yellow("Getting courses...\n")
	}

	url := fmt.Sprintf("https://%s%s?wstoken=%s&wsfunction=core_enrol_get_users_courses&userid=%s", types.Domain, types.Webservice, token, userID)
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

	courses := make([]types.Course, 0, len(names))
	name_def := ""
	for i, name := range names {
		// Localize the course name
		nameEs, nameEN := getCoursesNamesLanguages(name[1])
		if language == 1 {
			name_def = nameEs
		} else {
			name_def = nameEN
		}
		courses = append(courses, types.Course{Name: name_def, ID: ids[i][1]})
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
func SelectCourses(language int, coursesList []string, courses []types.Course) (bool, []string) {
	downloadAll := false
	if len(coursesList) != 0 && coursesList[0] == "all" {
		downloadAll = true
	} else if len(coursesList) == 0 {
		prompt := "Select the courses you want to download\n"
		if language == 1 {
			prompt = "Selecciona los cursos que quieres descargar\n"
		}
		listCoursesList := getCoursesNameByLanguage(courses)
		coursesList = checkboxesCourses(prompt, listCoursesList)
	}
	return downloadAll, coursesList
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
	survey.AskOne(prompt, &res, survey.WithKeepFilter(true))

	return res
}
