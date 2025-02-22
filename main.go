package main

import (
	"fmt"
	"log"

	c "github.com/Astrak00/AGDownloader/courses"
	download "github.com/Astrak00/AGDownloader/download"
	"github.com/Astrak00/AGDownloader/files"
	prog_args "github.com/Astrak00/AGDownloader/prog_args"
	token "github.com/Astrak00/AGDownloader/token"
	types "github.com/Astrak00/AGDownloader/types"
	u "github.com/Astrak00/AGDownloader/user"
	"github.com/fatih/color"
)

func main() {

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

	// use a full path instead of relative paths for the directory
	arguments.DirPath = files.GetFullPath(arguments.DirPath)

	fmt.Println("Downloading files to:", arguments.DirPath)

	// Obtain the user information by logging in with the token
	user, err := u.GetUserInfo(arguments.UserToken)
	for err != nil {
		user, err = u.GetUserInfo(arguments.UserToken)
		log.Default().Printf("Error getting user info: %v\nTrying again", err)
	}

	// Obtain the courses the user is enrolled in
	courses, err := c.GetCourses(arguments.UserToken, user.UserID, arguments.Language)
	if err != nil {
		log.Fatalf("Error getting courses: %v\n", err)
	}

	// Create an interactive list so the user can select the courses to download
	coursesList := c.SelectCoursesInteractive(arguments.Language, arguments.CoursesList, courses)

	// Create a channel to store the files and another for the errors that may occur when listing all the resources to download
	filesStoreChan := make(chan types.FileStore, len(courses)*100)
	errChan := make(chan error, len(courses))

	// List all the resources to downloaded and send them to the channel
	files.ListAllResources(coursesList, arguments.UserToken, arguments.DirPath, errChan, filesStoreChan)

	close(errChan)
	close(filesStoreChan)

	for err := range errChan {
		if err != nil {
			fmt.Println("Error:", err)
		}
	}

	// Download all the files in the channel
	download.DownloadFiles(filesStoreChan, arguments.MaxGoroutines, courses)

}
