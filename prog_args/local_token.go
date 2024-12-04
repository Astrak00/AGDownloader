package prog_args

import (
	"fmt"
	"log"
	"os"

	cookie "github.com/Astrak00/AGDownloader/cookies"
	"github.com/Astrak00/AGDownloader/types"
)

func ObtainingToken(arguments types.Prog_args) types.Prog_args {

	// Check if the token is stored in a local file to prevent unecessary request
	if _, err := os.Stat("token-file"); err == nil {
		data, err := os.ReadFile("token-file")
		if err != nil {
			log.Fatalf("Error reading file%v: %v\n", types.TokenDir, err)
		}
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
				fmt.Println("You must provide a cookie or a token to download the courses")
				fmt.Println(cookie.CookieText)
				fmt.Println()
				fmt.Println(cookie.ObtainCookieText)
				arguments.UserToken = cookie.GetTokenFromCookie(arguments)
			}
		}
	} else {
		arguments = AskForToken(arguments)
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
}
