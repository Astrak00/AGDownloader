package prog_args

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
	"regexp"
	"strings"
)

func ParseFlags() (int, string, string, int, []string) {
	// Definition of the flags used in this program
	languageStr := pflag.String("l", "ES", "Choose your language: ES: Español, EN:English")
	token := pflag.String("token", "", "Aula Global user security token 'aulaglobalmovil'")
	dir := pflag.String("dir", "", "Directory where you want to save the files")
	cores := pflag.Int("p", -1, "Cores to be used while downloading")
	var courses []string
	pflag.StringSliceVar(&courses, "courses", []string{}, "Ids or names of the courses to be downloaded, enclosed in \", separated by spaces. \n\"all\" downloads all courses")

	pflag.Parse()

	var language int
	if *languageStr == "ES" {
		language = 1
	} else {
		language = 2
	}

	// Attribution of the program creator
	if language == 1 {
		color.Cyan("Programa creado por Astrak00: github.com/Astrak00/AGDownloader/ \n" +
			"para descargar archivos de Aula Global en la UC3M\n")
	} else {
		color.Cyan("Program created by Astrak00: github.com/Astrak00/AGDownloader/ \n" +
			"to download files from Aula Global at UC3M\n")
	}

	// If the token or the directory are not given, prompt the user to introduce them
	if *dir == "" {
		*dir = promptForDir(language)
	}

	if *token == "" {
		*token = PromptForToken(language)
	}

	// If some courses are given, replace the commas with spaces and split the string
	if len(courses) == 1 && courses[0] != "" {
		courses[0] = strings.ReplaceAll(strings.ToLower(courses[0]), ",", " ")
		courses = strings.Split(courses[0], " ")
	}

	return language, *token, *dir, *cores, courses
}

// PromptForToken Prompt the user to introduce the token if it is not given
// Match the token with the regular expression to check if it is correct
// Correctness means that the token is at least 20 characters long and only contains letters and numbers
func PromptForToken(language int) string {
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

// Prompt the user to introduce the directory if it is not given
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
