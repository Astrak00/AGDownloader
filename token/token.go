package token

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/Astrak00/AGDownloader/cookies"
	"github.com/Astrak00/AGDownloader/types"

	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/all"
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

	fmt.Println("Please, open your browser and log in to Aula Global at UC3M")
	fmt.Println("Then, press enter to continue. If you have already logged in, press enter to continue")
	fmt.Scanln()

	if runtime.GOOS == "darwin" {
		fmt.Println("You will now be asked to enter your password to access the cookies. We need this to decrypt the cookie.")
	}

	// use kooky to get the token from the browser cookies
	var cookieResult string
	browserCookies := kooky.ReadCookies(kooky.DomainHasSuffix(`uc3m.es`), kooky.Name(`MoodleSessionag`))
	for _, cookie := range browserCookies {
		if len(cookie.Value) >= 26 {
			cookieResult = cookie.Value[len(cookie.Value)-26:]
		}
	}

	// get token from cookie
	if cookieResult == "" {
		cookieResult = cookies.AskForCookie()
	}
	token := cookies.CookieToToken(cookieResult)

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
