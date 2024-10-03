package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"

	c "github.com/Astrak00/AGDownloader/courses"
	download "github.com/Astrak00/AGDownloader/download"
	files "github.com/Astrak00/AGDownloader/files"
	types "github.com/Astrak00/AGDownloader/types"
	u "github.com/Astrak00/AGDownloader/user"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

func main() {
	language, userToken, dirPath, maxGoroutines, coursesList := parseFlags()

	// Obtain the user information by loggin in with the token
	user, err := u.GetUserInfo(userToken)
	for err != nil {
		user, err = u.GetUserInfo(promptForToken(language))
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
	var wg sync.WaitGroup
	files.ListAllResourcess(downloadAll, courses, userToken, dirPath, errChan, filesStoreChan, coursesList, &wg)

	wg.Wait()
	close(errChan)
	close(filesStoreChan)

	for err := range errChan {
		if err != nil {
			fmt.Println("Error:", err)
		}
	}

	download.DownloadFiles(filesStoreChan, maxGoroutines, language)

	if language == 1 {
		color.Green("Descarga completada\n")
	} else {
		color.Green("Download completed\n")
	}
}

func parseFlags() (int, string, string, int, []string) {
	language := pflag.Int("l", 0, "Choose your language: 1: Español, 2:English")
	token := pflag.String("token", "", "Aula Global user security token 'aulaglobalmovil'")
	dir := pflag.String("dir", "", "Directory where you want to save the files")
	cores := pflag.Int("p", 4, "Cores to be used while downloading")

	var courses []string
	pflag.StringSliceVar(&courses, "courses", []string{}, "Ids or names of the courses to be downloaded, enclosed in \", separated by spaces. \n\"all\" downloads all courses")
	pflag.Parse()

	if *language == 1 {
		color.Cyan("Programa creado por Astrak00: github.com/Astrak00/AGDownloader/ \n" +
			"para descargar archivos de Aula Global en la UC3M\n")
	} else {
		color.Cyan("Program created by Astrak00: github.com/Astrak00/AGDownloader/ \n" +
			"to download files from Aula Global at UC3M\n")
	}

	if *language == 0 {
		fmt.Println("Introduce tu idioma: 1: Español, 2:English")
		fmt.Scanf("%d", language)
	}

	if *dir == "" {
		*dir = promptForDir(*language)
	}

	if *token == "" {
		*token = promptForToken(*language)
	}

	// If some courses are given, replace the commas with spaces and split the string
	if len(courses) == 1 && courses[0] != "" {
		courses[0] = strings.ReplaceAll(strings.ToLower(courses[0]), ",", " ")
		courses = strings.Split(courses[0], " ")
	}

	return *language, *token, *dir, *cores, courses
}

func promptForToken(language int) string {
	var token string
	for {
		if language == 1 {
			fmt.Println("Introduzca el token de seguridad de su usuario de Aula Global 'aulaglobalmovil':")
		} else {
			fmt.Println("Introduce your Aula Global user security token 'aulaglobalmovil':")
		}
		fmt.Scanf("%s", &token)

		if token != "" && regexp.MustCompile(`[a-zA-Z0-9]{20,}`).MatchString(token) && len(token) > 20 {
			return token
		}

		if language == 1 {
			color.Red("El token introducido no parece estar correcto. Inténtelo de nuevo.")
		} else {
			color.Red("The given token does not seem to be right. Please try again.")
		}
	}
}

func promptForDir(language int) string {
	var dir string
	for {
		if language == 1 {
			fmt.Println("Introduzca la ruta donde quiere guardar los archivos:")
		} else {
			fmt.Println("Enter the path where you want to save the files:")
		}
		fmt.Scanf("%s", &dir)

		if dir != "" {
			return dir
		}

		if language == 1 {
			color.Red("La ruta introducida no parece estar correcta. Inténtelo de nuevo.")
		} else {
			color.Red("The given path does not seem to be right. Please try again.")
		}
	}
}
