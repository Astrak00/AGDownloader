package main

import (
	"fmt"
	c "github.com/Astrak00/AGDownloader/courses"
	download "github.com/Astrak00/AGDownloader/download"
	files "github.com/Astrak00/AGDownloader/files"
	prog_args "github.com/Astrak00/AGDownloader/prog_args"
	types "github.com/Astrak00/AGDownloader/types"
	u "github.com/Astrak00/AGDownloader/user"
	"github.com/fatih/color"
	"log"
)

func main() {
	// Parse the flags to get the language, userToken, dirPath, maxGoroutines and coursesList

	language, userToken, dirPath, maxGoroutines, coursesList := prog_args.ParseFlags()

	// Obtain the user information by loggin in with the token
	user, err := u.GetUserInfo(userToken)
	for err != nil {
		user, err = u.GetUserInfo(prog_args.PromptForToken(language))
	}

	// Obtain the courses of the user
	courses, err := c.GetCourses(userToken, user.UserID, language)
	if err != nil {
		log.Fatalf("Error getting courses: %v\n", err)
	}

	// Create an interactive list so the user can select the courses to download
	downloadAll, coursesList := c.SelectCourses(language, coursesList, courses)

	// Create a channel to store the files and another for the errors that may occur when listing all the resources to download
	filesStoreChan := make(chan types.FileStore, len(courses)*20)
	errChan := make(chan error, len(courses))

	// Creating a chanel to store the files that wull be downloaded
	files.ListAllResourcess(downloadAll, courses, userToken, dirPath, errChan, filesStoreChan, coursesList)

	close(errChan)
	close(filesStoreChan)

	for err := range errChan {
		if err != nil {
			fmt.Println("Error:", err)
		}
	}

	// Download all the files in the channel
	if maxGoroutines == -1 {
		maxGoroutines = len(filesStoreChan)
	}
	download.DownloadFiles(filesStoreChan, maxGoroutines)

	if language == 1 {
		color.Green("Descarga completada\n")
	} else {
		color.Green("Download completed\n")
	}
}
