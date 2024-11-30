package prog_args

import (
	"fmt"
	"log"
	"strconv"

	types "github.com/Astrak00/AGDownloader/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

func ParseFlags() types.Prog_args {
	// Definition of the flags used in this program
	languageStr := pflag.String("l", "EN", "Choose your language: ES: Español, EN:English")
	token := pflag.String("token", "", "Aula Global user security token 'aulaglobalmovil'")
	dir := pflag.String("dir", "", "Directory where you want to save the files")
	cores := pflag.Int("p", -1, "Cores to be used while downloading")
	var courses []string
	pflag.StringSliceVar(&courses, "courses", []string{}, "Ids or names of the courses to be downloaded, enclosed in \", separated by spaces. \n\"all\" downloads all courses")

	pflag.Parse()

	var language int
	switch *languageStr {
	case "ES":
		language = 1
	default:
		language = 2
	}

	return types.Prog_args{
		Language:      language,
		UserToken:     *token,
		DirPath:       *dir,
		MaxGoroutines: *cores,
		CoursesList:   courses,
	}
}

// Knowledge_token Ask the user if they know their token
func AskForToken(arguments types.Prog_args) types.Prog_args {

	// Attribution of the program creator
	if arguments.Language == 1 {
		color.Cyan("Programa creado por Astrak00 para descargar archivos de Aula Global en la UC3M\n")
	} else {
		color.Cyan("Program created by Astrak00 to download files from Aula Global at UC3M\n")
	}

	p := tea.NewProgram(initialModel(&arguments.DirPath, &arguments.UserToken, arguments.MaxGoroutines))

	finalModel, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	tokenObtained := finalModel.(model).inputs[tokenIota].Value()
	dirObtained := finalModel.(model).inputs[dirIota].Value()
	coresObtained, err := strconv.Atoi(finalModel.(model).inputs[corIota].Value())
	if err != nil {
		log.Fatal(err)
	}

	return types.Prog_args{
		Language:      arguments.Language,
		UserToken:     tokenObtained,
		DirPath:       dirObtained,
		MaxGoroutines: coresObtained,
		CoursesList:   arguments.CoursesList,
	}
}

// PromptForToken Prompt the user to introduce the token if it is not given
// Match the token with the regular expression to check if it is correct
// Correctness means that the token is at least 20 characters long and only contains letters and numbers
func PromptForToken(language int) string {
	var token string
	for {
		if language == 1 {
			color.Yellow("Ha habido un error con el token, por favor, introdúcelo de nuevo:")
		} else {
			color.Yellow("There has been an error with the token, please input it again:")
		}
		_, err := fmt.Scanf("%s", &token)
		if err != nil {
			return ""
		}

		if tokenValidator(token) == nil {
			return token
		}

		if language == 1 {
			color.Red("El token introducido no parece estar correcto. Inténtelo de nuevo.")
		} else {
			color.Red("The given token does not seem to be right. Please try again.")
		}
	}
}