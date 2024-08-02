package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	courseParser "github.com/100472175/AGDownloader/functions"

	"github.com/fatih/color"
)

// This function processes the course and downloads all the files available
func processCourse(wg *sync.WaitGroup, token string, course courseParser.Course, directory string) {
	defer wg.Done()
	fmt.Printf("Course: %s\n", course.Name)
	modules, err := courseParser.GetCourseContent(token, course.ID)
	if err != nil {
		color.Red("Error: %s", err)
		return
	}
	if err := courseParser.SaveFiles(token, course.Name, modules, directory); err != nil {
		color.Red("Error: %s", err)
	}
}

func main() {
	color.Blue("Welcome to AGDownloader. To use this cli tool, you will need the \"aulaglobalmovil\" key. Use the -help flag to see the usage.")

	// Set the flags for the cli tool: the token, which will be used to authenticate, the directory where the files will be saved and the help flag.
	token_ptr := flag.String("token", "00", "aulaglobalmobile token used to authenticate")
	directory_ptr := flag.String("directory", "downloadedCourseFiles", "Directory where the files will be saved")
	help_flag_ptr := flag.Bool("help", false, "Show help")
	flag.Parse()

	token := *token_ptr
	directory := *directory_ptr
	if *help_flag_ptr || token == "00" {
		fmt.Println("Usage: AGDownloader --token <token> --directory <directory>")
		color.Yellow("To obtain the token, you need to log in to Aula Global and go to the preferences panel. There, select security keys and copy or regenerate if expired, the \"aulaglobalmovil\" key.")
		return
	}

	// Get the user ID
	userID, _, err := courseParser.GetUserInfo(token)
	if err != nil {
		color.Red("Error: %s", err)
		return
	}
	color.Yellow("Obtaining courses...")

	// Parse the courses available
	courses, err := courseParser.ParseXmlCourses(token, userID)
	if err != nil {
		color.Red("Error: %s", err)
		return
	}

	startTime := time.Now()
	wg := sync.WaitGroup{}
	for _, course := range courses {
		wg.Add(1)
		go processCourse(&wg, token, course, directory)
	}
	wg.Wait()
	elapsedTime := time.Since(startTime)
	color.Green("All files downloaded in %s", elapsedTime)
}
