package prog_args

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	types "github.com/Astrak00/AGDownloader/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/pflag"
)

func tokenValidator(s string) error {
	// Token should be a string of more than 20 characters, that matches the regular expression
	if s != "" && len(s) > 20 && regexp.MustCompile(`[a-zA-Z0-9]{20,}`).MatchString(s) {
		return nil
	}
	return fmt.Errorf("token is invalid")
}

/*
ParseCLIArgs parses the command-line arguments and returns a ProgramArgs struct.
It defines and processes the following flags:

-l: Language of the course names, either "ES" (Español) or "EN" (English). Default is "ES".

--token: Aula Global user security token 'aulaglobalmovil'.

--dir: Directory where the files will be saved.

-p: Number of cores to be used while downloading. Default is 0.

--fast: If set, MaxGoroutines will be set to the number of files for fastest downloading.

--courses: A list of course IDs or names to be downloaded, enclosed in quotes and separated by spaces. "all" downloads all courses.
It validates the token and adjusts the number of cores if the fast flag is set.

Returns a ProgramArgs struct containing the parsed values.
*/
func ParseCLIArgs() types.ProgramArgs {
	languageStr := pflag.String("l", "ES", "Language of the course names: ES (Español) or EN (English)")
	token := pflag.String("token", "", "Aula Global user security token 'aulaglobalmovil'")
	dir := pflag.String("dir", "", "Directory where you want to save the files")
	cores := pflag.Int("p", 0, "Number of cores to be used while downloading")
	fast := pflag.Bool("fast", false, "Set MaxGoroutines to the number of files for fastest downloading")
	webUI := pflag.Bool("web", false, "Select the courses using the web interface")
	participantsList := pflag.Bool("pList", false, "Select weather you want to also download the participants list of the course")
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

	if *token != "" {
		if err := tokenValidator(*token); err != nil {
			log.Fatalf("Error getting courses: %v\n", err)
		}
	}

	if *fast {
		*cores = -1
	}

	return types.ProgramArgs{
		Language:         language,
		UserToken:        *token,
		DirPath:          *dir,
		MaxGoroutines:    *cores,
		CoursesList:      courses,
		WebUI:            *webUI,
		ParticipantsList: *participantsList,
	}
}

// Ask the user in an interactive way for the missing arguments.
// The user will be prompted for the directory path and the number of cores to use.
// Returns a ProgramArgs struct with the obtained values.
func PromptMissingArgs(arguments types.ProgramArgs) types.ProgramArgs {

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

	return types.ProgramArgs{
		Language:         arguments.Language,
		UserToken:        arguments.UserToken,
		DirPath:          dirObtained,
		MaxGoroutines:    coresObtained,
		CoursesList:      arguments.CoursesList,
		WebUI:            arguments.WebUI,
		ParticipantsList: arguments.ParticipantsList,
	}
}
