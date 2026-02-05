package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	c "github.com/Astrak00/AGDownloader/courses"
	download "github.com/Astrak00/AGDownloader/download"
	errorlog "github.com/Astrak00/AGDownloader/errorlog"
	"github.com/Astrak00/AGDownloader/files"
	prog_args "github.com/Astrak00/AGDownloader/prog_args"
	token "github.com/Astrak00/AGDownloader/token"
	types "github.com/Astrak00/AGDownloader/types"
	u "github.com/Astrak00/AGDownloader/user"
	webui "github.com/Astrak00/AGDownloader/webUI"
	"github.com/fatih/color"
)

func main() {
	// Set up global signal handling
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		fmt.Println("\nReceived interrupt signal, exiting...")
		os.Exit(0)
	}()

	// Parse the flags to get the language, user token, the path to save the downloaded files, maxGoroutines and courses list to download
	arguments := prog_args.ParseCLIArgs()

	// Attribution of the program creator
	color.Cyan("Program created by Astrak00 to download files from Aula Global at UC3M\n")

	// In case the user has not provided a token though the cli, we try to obtain it from a file or ask the user for it
	if arguments.UserToken == "" {
		arguments.UserToken = token.ObtainToken()
	}

	// If there are missing arguments, we prompt the user for them
	if !arguments.CheckAllAsigned() {
		arguments = prog_args.PromptMissingArgs(arguments)
	}

	// Initialize error logger
	errLogger, err := errorlog.New(arguments.DirPath)
	if err != nil {
		log.Printf("Warning: Failed to initialize error logger: %v\n", err)
		log.Println("Continuing without error logging...")
		errLogger = nil
	} else {
		defer func() {
			if errLogger != nil {
				errCount := errLogger.GetErrorCount()
				if errCount > 0 {
					color.Yellow("\nTotal errors logged: %d\n", errCount)
					color.Yellow("Error log saved to: %s\n", errLogger.GetLogFilePath())
				}
				errLogger.Close()
			}
		}()
		color.Green("Error logging initialized: %s\n", errLogger.GetLogFilePath())
	}

	// Obtain the user information by logging in with the token
	user, err := u.GetUserInfo(arguments.UserToken)
	retriesCounter := 0
	for err != nil && retriesCounter < 3 {
		user, err = u.GetUserInfo(arguments.UserToken)
		retriesCounter++
		log.Default().Printf("Error getting user info: %v\nTrying again. Attempt %d/3", err, retriesCounter)
		if retriesCounter == 3 {
			log.Fatalf("Error getting user info after 3 attempts")
		}
	}

	// Obtain the courses the user is enrolled in
	var courses types.Courses
	if arguments.Timeline {
		// Use timeline API to get all courses (current, past, and future)
		courses, err = c.GetCoursesByTimeline(arguments.UserToken, arguments.Language)
	} else {
		// Use standard API with userID
		courses, err = c.GetCourses(arguments.UserToken, user.UserID, arguments.Language)
	}
	if err != nil {
		log.Fatalf("Error getting courses: %v\n", err)
	}

	var coursesList []types.Course
	if arguments.WebUI {
		coursesList = webui.ShowCourseWeb(courses)
	} else {
		coursesList = c.SelectCoursesInteractive(arguments.Language, arguments.CoursesList, courses)
	}
	// Create an interactive list so the user can select the courses to download

	// Create a channel to store the files and another for the errors that may occur when listing all the resources to download
	filesStoreChan := make(chan types.FileStore, len(courses)*100)
	errChan := make(chan error, len(courses))

	// List all the resources to downloaded and send them to the channel
	files.ListAllResources(coursesList, arguments.UserToken, arguments.DirPath, errChan, filesStoreChan, errLogger)

	close(errChan)
	close(filesStoreChan)

	for err := range errChan {
		if err != nil {
			fmt.Println("Error:", err)
		}
	}

	// Download all the files in the channel
	download.DownloadFiles(filesStoreChan, arguments.MaxGoroutines, coursesList, errLogger)

}
