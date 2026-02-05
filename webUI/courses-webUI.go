package webui

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/Astrak00/AGDownloader/types"
)

var courseList []types.Course

func ShowCourseWeb(courses []types.Course) []types.Course {
	courseList = courses

	// Handle the main page (the form) and the submission endpoint.
	http.HandleFunc("/", formHandler)
	// http.HandleFunc("/submit", submitHandler)

	// fmt.Println("Server running on http://localhost:8888")
	// Open the default browser with the URL.
	err := exec.Command("open", "http://localhost:8888").Start()
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Start the server and wait for the form submission.
	selectedCourses := make(chan []string)
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		// Parse the form values.
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		// Get the list of selected courses.
		selected := r.Form["courses"]
		selectedCourses <- selected

		w.Header().Set("Content-Type", "text/html")

		correctResponseByte := []byte(correctResponseHTML)
		if _, err := w.Write(correctResponseByte); err != nil {
			http.Error(w, "Error writing response", http.StatusInternalServerError)
			return
		}

	})

	go func() {
		if err := http.ListenAndServe(":8888", nil); err != nil {
			fmt.Println("Server error:", err)
		}
	}()

	// Wait for the selected courses.
	selected := <-selectedCourses
	var selectedCourseList []types.Course
	for _, courseName := range selected {
		for _, course := range courseList {
			if course.Name == courseName {
				selectedCourseList = append(selectedCourseList, course)
				break
			}
		}
	}
	return selectedCourseList
}

// formHandler serves the HTML form with checkboxes for each letter.
func formHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	courseSelectorHTMLByte := []byte(courseSelectorHTMLStart)
	if _, err := w.Write(courseSelectorHTMLByte); err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}

	// Loop through the letters A-Z to create a checkbox for each.
	for _, course := range courseList {
		fmt.Fprintf(w, ` <label class="course-option" for="course%s">
                <input type="checkbox" id="course%s" name="courses" value="%s" class="hidden-checkbox">
                <div class="checkmark">✓</div>
                <div class="course-title">%s</div>
            </label>`, course.Name, course.ID, course.Name, course.Name)
	}

	courseSelectorHTMLByte = []byte(courseSelectorHTMLEnd)
	if _, err := w.Write(courseSelectorHTMLByte); err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}
