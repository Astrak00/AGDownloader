package token

import (
	"fmt"
	"log"
	"os"

	"github.com/Astrak00/AGDownloader/cookies"
	"github.com/Astrak00/AGDownloader/types"
	webui "github.com/Astrak00/AGDownloader/webUI"
)

// ObtainToken gets the token from the saved file from a previous execution or asks the user for it
// and saves it to a file.
// Returns the token.
func ObtainToken() string {

	// Check if the token is stored in a local file to prevent unecessary request
	if _, err := os.Stat(types.TokenDir); err == nil {
		data, err := os.ReadFile(types.TokenDir)
		if err != nil {
			log.Fatalf("Error reading file%v: %v\n", types.TokenDir, err)
		}
		//fmt.Println("Token token loaded from", types.TokenDir)
		return string(data)
	}

	// get token from cookie using web popup
	fmt.Println("Opening browser to obtain cookie...")
	cookie := webui.AskForCookieWeb()
	if cookie == "" {
		cookie = cookies.AskForCookie()
	}
	token := cookies.CookieToToken(cookie)

	saveToken(token)

	return token
}

// saveToken saves the token to a file names types.TokenDir (aulaglobal-token)
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
