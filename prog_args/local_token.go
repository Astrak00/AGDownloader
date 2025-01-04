package prog_args

import (
	"fmt"
	"log"
	"os"
	"time"

	cookie "github.com/Astrak00/AGDownloader/cookies"
	"github.com/Astrak00/AGDownloader/types"
	"github.com/fatih/color"
)

func ObtainingToken(arguments types.Prog_args) types.Prog_args {

	// Check if the token is stored in a local file to prevent unecessary request
	if _, err := os.Stat(types.TokenDir); err == nil {
		data, err := os.ReadFile(types.TokenDir)
		if err != nil {
			log.Fatalf("Error reading file%v: %v\n", types.TokenDir, err)
		}
		//fmt.Println("Token token loaded from", types.TokenDir)
		arguments.UserToken = string(data)
	}

	// Ask the user if they have the token from another place
	if arguments.UserToken == "" {
		if Knowledge_element("token") {
			arguments = AskForToken(arguments)
		} else {
			// Ask for the cookie
			if Knowledge_element("auth cookie") {
				arguments.UserToken = cookie.GetTokenFromCookie(arguments)
			} else {
				color.Yellow(cookie.CookieText)
				color.Magenta(cookie.ObtainCookieText)
				time.Sleep(2 * time.Second)
				arguments.UserToken = cookie.GetTokenFromCookie(arguments)
			}
		}
	} else {
		if arguments.UserToken == "" || arguments.DirPath == "" || arguments.MaxGoroutines == 0 {
			arguments = AskForToken(arguments)
		}
	}

	saveToken(arguments.UserToken)
	return arguments
}

func saveToken(token string) {
	if token == "" {
		return
	}
	// We save the token to a file to be able to read it in future executions
	err := os.WriteFile(types.TokenDir, []byte(token), 0644)
	if err != nil {
		log.Fatal("Error saving the token to a file", err)
	}
	fmt.Println("Token saved to", types.TokenDir)
}
