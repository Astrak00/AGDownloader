package main

import (
	"fmt"
	"log"

	c "github.com/Astrak00/AGDownloader/courses"
	download "github.com/Astrak00/AGDownloader/download"
	files "github.com/Astrak00/AGDownloader/files"
	prog_args "github.com/Astrak00/AGDownloader/prog_args"
	types "github.com/Astrak00/AGDownloader/types"
	u "github.com/Astrak00/AGDownloader/user"
	"github.com/fatih/color"
)

func main() {
	// Parse the flags to get the language, userToken, dirPath, maxGoroutines and coursesList

	arguments := prog_args.ParseFlags()

	arguments = prog_args.ObtainingToken(arguments)

	// Obtain the user information by loggin in with the token
	user, err := u.GetUserInfo(arguments.UserToken)
	for err != nil {
		user, err = u.GetUserInfo(arguments.UserToken)
	}

	// Obtain the courses of the user
	courses, err := c.GetCourses(arguments.UserToken, user.UserID, arguments.Language)
	if err != nil {
		log.Fatalf("Error getting courses: %v\n", err)
	}

	// Create an interactive list so the user can select the courses to download
	downloadAll, coursesList := c.SelectCourses(arguments.Language, arguments.CoursesList, courses)

	// Create a channel to store the files and another for the errors that may occur when listing all the resources to download
	filesStoreChan := make(chan types.FileStore, len(courses)*20)
	errChan := make(chan error, len(courses))

	// Creating a chanel to store the files that wull be downloaded
	files.ListAllResourcess(downloadAll, courses, arguments.UserToken, arguments.DirPath, errChan, filesStoreChan, coursesList)

	close(errChan)
	close(filesStoreChan)

	for err := range errChan {
		if err != nil {
			fmt.Println("Error:", err)
		}
	}

	// Download all the files in the channel
	if arguments.MaxGoroutines == -1 {
		arguments.MaxGoroutines = len(filesStoreChan)
	}
	download.DownloadFiles(filesStoreChan, arguments.MaxGoroutines)

	if arguments.Language == 1 {
		color.Green("Descarga completada\n")
	} else {
		color.Green("Download completed\n")
	}
}
