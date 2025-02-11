package prog_args

import (
	"log"
	"strconv"
	"fmt"
	"regexp"

	types "github.com/Astrak00/AGDownloader/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/pflag"
)


func tokenValidator(s string) error {
	// Token should be a string of 22 characters, that matches the regular expression
	if s != "" && regexp.MustCompile(`[a-zA-Z0-9]{20,}`).MatchString(s) && len(s) > 20 {
		return nil
	}
	return fmt.Errorf("token is invalid")
}

func ParseFlags() types.Prog_args {
	// Definition of the flags used in this program
	languageStr := pflag.String("l", "ES", "Language of the course names: ES (Español) or EN (English)")
	token := pflag.String("token", "", "Aula Global user security token 'aulaglobalmovil'")
	dir := pflag.String("dir", "", "Directory where you want to save the files")
	cores := pflag.Int("p", 0, "Number of cores to be used while downloading")
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

	// validate token
	if *token != "" {
		if err := tokenValidator(*token); err != nil {
			log.Fatalf("Error getting courses: %v\n", err)
		}
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

// get the program config
func GetConfig(arguments types.Prog_args) types.Prog_args {

	p := tea.NewProgram(initialModel(&arguments.DirPath, arguments.MaxGoroutines))

	finalModel, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	dirObtained := finalModel.(model).inputs[dirIota].Value()
	if dirObtained == "" {
		dirObtained = "."
	}
	coresObtained, err := strconv.Atoi(finalModel.(model).inputs[corIota].Value())
	if err != nil {
		log.Fatal(err)
	}

	return types.Prog_args{
		Language:      arguments.Language,
		UserToken:     arguments.UserToken,
		DirPath:       dirObtained,
		MaxGoroutines: coresObtained,
		CoursesList:   arguments.CoursesList,
	}
}
