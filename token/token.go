package token

import (
	"fmt"
	"log"
	"os"

	"github.com/Astrak00/AGDownloader/cookies"
	"github.com/Astrak00/AGDownloader/types"
)

func ObtaininToken() string {

	// Check if the token is stored in a local file to prevent unecessary request
	if _, err := os.Stat(types.TokenDir); err == nil {
		data, err := os.ReadFile(types.TokenDir)
		if err != nil {
			log.Fatalf("Error reading file%v: %v\n", types.TokenDir, err)
		}
		//fmt.Println("Token token loaded from", types.TokenDir)
		return string(data)
	}

	// get token from cookie
	cookie := cookies.AskForCookie()
	token := cookies.CookieToToken(cookie)

	saveToken(token)

	return token
}

func saveToken(token string) {
	if token == "" {
		return
	}

	// TODO: prompt user for save path, etc.

	// We save the token to a file to be able to read it in future executions
	err := os.WriteFile(types.TokenDir, []byte(token), 0644)
	if err != nil {
		log.Fatal("Error saving the token to a file", err)
	}
	fmt.Println("Token saved to", types.TokenDir)
}
