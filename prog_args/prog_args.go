package prog_args

import (
	"log"
	"strconv"

	types "github.com/Astrak00/AGDownloader/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/pflag"
)

func ParseFlags() types.Prog_args {
	// Definition of the flags used in this program
	languageStr := pflag.String("l", "EN", "Choose your language: ES: Español, EN:English")
	token := pflag.String("token", "", "Aula Global user security token 'aulaglobalmovil'")
	dir := pflag.String("dir", "", "Directory where you want to save the files")
	cores := pflag.Int("p", 0, "Cores to be used while downloading")
	fast := pflag.Bool("fast", false, "Set MaxGoroutines to the number of files for fastest downloading")
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

	if *fast {
		*cores = -1
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

	p := tea.NewProgram(initialModel(&arguments.DirPath, &arguments.UserToken, arguments.MaxGoroutines))

	finalModel, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	tokenObtained := finalModel.(model).inputs[tokenIota].Value()
	dirObtained := finalModel.(model).inputs[dirIota].Value()
	if dirObtained == "" {
		dirObtained = "downloaded_files"
	}
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
